package connect

import (
	"net"
	"sync"
)

// ConnectionPool manages a pool of network connections.
type ConnectionPool struct {
	maxIdle     int
	connections chan net.Conn
	dial        func() (net.Conn, error)
	mutex       sync.Mutex
	validate    func(conn net.Conn) bool
	closed      bool
	once        sync.Once
}

// NewConnectionPool creates a new connection pool.
func NewConnectionPool(maxIdle int, dial func() (net.Conn, error)) *ConnectionPool {
	return &ConnectionPool{
		maxIdle:     maxIdle,
		connections: make(chan net.Conn, maxIdle),
		dial:        dial,
	}
}

// Get retrieves a connection from the pool.
func (p *ConnectionPool) Get() (net.Conn, error) {
	for {
		select {
		case conn := <-p.connections:
			if p.validate == nil || p.validate(conn) {
				return conn, nil
			}
			conn.Close() // invalid connection
		default:
			return p.dial()
		}
	}
}

// Put returns a connection to the pool.
func (p *ConnectionPool) Put(conn net.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	select {
	case p.connections <- conn:
	default:
		conn.Close()
	}
}

// Close closes all connections in the pool.
func (p *ConnectionPool) Close() {
	p.once.Do(func() {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		p.closed = true
		close(p.connections)
		for conn := range p.connections {
			conn.Close()
		}
	})

}
