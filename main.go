package main

import (
	"fmt"
	"github.com/kenelite/gedis/internal/connect"
	"github.com/kenelite/gedis/internal/handle"
	"github.com/kenelite/gedis/internal/response"
	"github.com/kenelite/gedis/internal/storage"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new response
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := storage.NewAof("./database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	// creat connection pool
	pool := connect.NewConnectionPool(5, func() (net.Conn, error) {
		return net.Dial("tcp", "localhost:6379")
	})
	defer pool.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, pool, aof)
	}
}

func handleConnection(conn net.Conn, pool *connect.ConnectionPool, aof *storage.Aof) {
	defer conn.Close()

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
			fmt.Println(err)
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

		handler, ok := handle.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(response.Value{Typ: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
