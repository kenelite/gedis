package main

import (
	"fmt"
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

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

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
