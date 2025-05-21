package main

import (
	"fmt"
	gedis "github.com/kenelite/gedis/gedis-go-sdk"
	"log"
)

func main() {
	client, err := gedis.NewClient("localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	_, _ = client.Auth("admin", "admin123")
	_, _ = client.Set("foo", "bar")
	val, _ := client.Get("foo")
	fmt.Println("GET foo:", val)

	_, _ = client.LPush("mylist", "one")
	_, _ = client.RPush("mylist", "two")
}
