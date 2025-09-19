package sse

type Client struct {
	id    string
	send  chan []byte
	close chan struct{}
}

func NewClient(id string) *Client {
	return &Client{
		id:    id,
		send:  make(chan []byte, 256),
		close: make(chan struct{}),
	}
}

func (c *Client) Send() chan []byte {
	return c.send
}

func (c *Client) Close() chan struct{} {
	return c.close
}
