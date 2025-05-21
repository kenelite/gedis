package main

import (
	"fmt"
	"os"

	"github.com/kenelite/gedis-cli/client"
)

func main() {
	addr := "localhost:6379" // or read from flags/env

	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	c, err := client.NewClient(addr)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		os.Exit(1)
	}

	c.Run()
}
