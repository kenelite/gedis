package gedis_go_sdk

type Client struct {
	conn *Conn
}

// NewClient connects to a Gedis server at addr (e.g. "localhost:6379")
func NewClient(addr string) (*Client, error) {
	conn, err := NewConn(addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Auth(user, password string) (string, error) {
	return c.conn.SendCommand("AUTH", user, password)
}

// SET key value
func (c *Client) Set(key, value string) (string, error) {
	return c.conn.SendCommand("SET", key, value)
}

// GET key
func (c *Client) Get(key string) (string, error) {
	return c.conn.SendCommand("GET", key)
}

// DEL key
func (c *Client) Del(key string) (string, error) {
	return c.conn.SendCommand("DEL", key)
}

// LPUSH key value
func (c *Client) LPush(key string, value string) (string, error) {
	return c.conn.SendCommand("LPUSH", key, value)
}

// RPUSH key value
func (c *Client) RPush(key string, value string) (string, error) {
	return c.conn.SendCommand("RPUSH", key, value)
}

// HSET key field value
func (c *Client) HSet(key, field, value string) (string, error) {
	return c.conn.SendCommand("HSET", key, field, value)
}

// HGET key field
func (c *Client) HGet(key, field string) (string, error) {
	return c.conn.SendCommand("HGET", key, field)
}
