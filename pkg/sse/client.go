package sse

type Client struct {
	send chan []byte
}

func NewClient() *Client {
	return &Client{
		send: make(chan []byte, 256),
	}
}

func (c *Client) Send() chan []byte {
	return c.send
}
