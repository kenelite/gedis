package gedis_go_sdk

import (
	"bufio"
	"net"
	"strings"
)

type Conn struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewConn(addr string) (*Conn, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn:   c,
		reader: bufio.NewReader(c),
		writer: bufio.NewWriter(c),
	}, nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

// SendCommand formats and sends a RESP command and reads the response
func (c *Conn) SendCommand(command string, args ...string) (string, error) {
	request := formatRespArray(command, args...)
	_, err := c.writer.WriteString(request)
	if err != nil {
		return "", err
	}
	err = c.writer.Flush()
	if err != nil {
		return "", err
	}

	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
