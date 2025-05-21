package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/kenelite/gedis-cli/resp"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

func (c *Client) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("gedis-cli connected. Type 'exit' to quit.")

	for {
		fmt.Print("gedis> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "exit" {
			break
		}
		if line == "" {
			continue
		}

		args := strings.Fields(line)
		if err := resp.EncodeArray(c.writer, args); err != nil {
			fmt.Println("ERR:", err)
			continue
		}
		c.writer.Flush()

		respVal, err := resp.Decode(c.reader)
		if err != nil {
			fmt.Println("ERR:", err)
			continue
		}
		fmt.Println(respVal)
	}
}
