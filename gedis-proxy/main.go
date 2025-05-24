package main

import (
	"github.com/kenelite/gedis/proxy/proxy"
	"log"
)

func main() {
	backends := []string{
		"localhost:6380",
	}

	balancer := proxy.NewBalancer(backends)

	server := &proxy.ProxyServer{
		ListenAddr: ":6379",
		Balancer:   balancer,
	}

	if err := server.Start(); err != nil {
		log.Fatal("Failed to start proxy:", err)
	}
}
