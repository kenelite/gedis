package proxy

import (
	"net"
	"sync"
)

type Balancer struct {
	backends []string
	current  int
	mu       sync.Mutex
}

func NewBalancer(backends []string) *Balancer {
	return &Balancer{backends: backends}
}

func (b *Balancer) Next() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	addr := b.backends[b.current]
	b.current = (b.current + 1) % len(b.backends)
	return addr
}

func (b *Balancer) GetConnection() (net.Conn, error) {
	for i := 0; i < len(b.backends); i++ {
		addr := b.Next()
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			return conn, nil
		}
	}
	return nil, ErrAllBackendsDown
}

var ErrAllBackendsDown = &net.OpError{
	Op:   "dial",
	Net:  "tcp",
	Addr: nil,
	Err:  errAllBackends,
}

var errAllBackends = &net.AddrError{
	Err:  "all backends are down",
	Addr: "",
}
