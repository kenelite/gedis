package main

import (
	"fmt"
	"github.com/kenelite/gedis/internal/auth"
	"github.com/kenelite/gedis/internal/config"
	"github.com/kenelite/gedis/internal/connect"
	"github.com/kenelite/gedis/internal/handle"
	"github.com/kenelite/gedis/internal/response"
	"github.com/kenelite/gedis/internal/storage"
	"net"
	"strconv"
	"strings"
)

func main() {

	// load user and passport
	err := auth.LoadUsersFromConfig("settings/users.conf")
	if err != nil {
		panic(err)
	}

	// load config
	cfg, err := config.Load("settings/server.conf")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	port := cfg.Get("server", "port")
	if port == "" {
		port = "6379" // default port
	}

	aofPath := cfg.Get("server", "aof_path")
	if aofPath == "" {
		aofPath = "settings/database.aof"
	}

	aofSyncIntervalStr := cfg.Get("server", "aof_sync_interval_sec")
	aofSyncInterval := 1 // default 1 second
	if aofSyncIntervalStr != "" {
		if val, err := strconv.Atoi(aofSyncIntervalStr); err == nil {
			aofSyncInterval = val
		}
	}

	// boot
	fmt.Println("Listening on port :6379")

	// Create a new response
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := storage.NewAofWithInterval(aofPath, aofSyncInterval)
	if err != nil {
		fmt.Println(err)
		return
	}
	aof.Load()
	defer aof.Close()

	// creat connection pool
	pool := connect.NewConnectionPool(5, func() (net.Conn, error) {
		return net.Dial("tcp", "localhost:6379")
	})
	defer pool.Close()

	handle.StartKeyExpirationLoop()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, pool, aof)
	}
}

func handleConnection(conn net.Conn, pool *connect.ConnectionPool, aof *storage.Aof) {
	defer conn.Close()

	authenticated := false
	pooledConn, err := pool.Get()
	if err != nil {
		fmt.Println("Error getting connection from pool:", err)
		return
	}
	defer pool.Put(pooledConn)

	_, err = pooledConn.Write([]byte("Hello from pool!\n"))
	if err != nil {
		fmt.Println("Error writing to pooled connection:", err)
		return
	}

	for {
		resp := response.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println("RESP Read error:", err)
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := response.NewWriter(conn)

		// Special handling for AUTH command
		if command == "AUTH" {
			if len(args) != 2 {
				writer.Write(response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'AUTH'"})
				continue
			}
			username := args[0].Bulk
			password := args[1].Bulk
			if auth.CheckUser(username, password) {
				authenticated = true
				writer.Write(response.Value{Typ: "string", Str: "OK"})
			} else {
				writer.Write(response.Value{Typ: "error", Str: "ERR invalid username or password"})
			}
			continue
		}

		// Require authentication for all commands except AUTH
		if !authenticated {
			writer.Write(response.Value{Typ: "error", Str: "NOAUTH Authentication required"})
			continue
		}

		handler, ok := handle.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(response.Value{Typ: "error", Str: "ERR unknown command: " + command})
			continue
		}

		result := handler(args)
		if result.Typ != "error" && (command == "SET" || command == "HSET" || command == "LPUSH" || command == "RPUSH") {
			aof.Write(value)

			// If SET with EX, write an EXPIRE command
			if command == "SET" && len(value.Array) >= 4 && strings.ToUpper(value.Array[2].Bulk) == "EX" {
				key := value.Array[0].Bulk
				ttl := value.Array[3].Bulk

				expireCmd := response.Value{
					Typ: "array",
					Array: []response.Value{
						{Typ: "bulk", Bulk: "EXPIRE"},
						{Typ: "bulk", Bulk: key},
						{Typ: "bulk", Bulk: ttl},
					},
				}
				aof.Write(expireCmd)
			}
		}

		writer.Write(result)
	}
}
