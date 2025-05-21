package proxy

import (
	"io"
	"log"
	"net"
)

type ProxyServer struct {
	ListenAddr string
	Balancer   *Balancer
}

func (p *ProxyServer) Start() error {
	ln, err := net.Listen("tcp", p.ListenAddr)
	if err != nil {
		return err
	}
	log.Println("Proxy listening on", p.ListenAddr)

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept:", err)
			continue
		}
		go p.handleClient(clientConn)
	}
}

func (p *ProxyServer) handleClient(client net.Conn) {
	defer client.Close()

	backendConn, err := p.Balancer.GetConnection()
	if err != nil {
		log.Println("All backends are down.")
		client.Write([]byte("-ERR all backends are down\r\n"))
		return
	}
	defer backendConn.Close()

	go io.Copy(backendConn, client) // Client -> Backend
	io.Copy(client, backendConn)    // Backend -> Client
}
