package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &MyTelnetClient{addr: address, timeout: timeout, in: in, out: out}
}

type MyTelnetClient struct {
	addr    string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (c *MyTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return fmt.Errorf("connect:%w", err)
	}
	c.conn = conn
	return nil
}

func (c *MyTelnetClient) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("close:%w", err)
		}
	}
	return nil
}

func (c *MyTelnetClient) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("send:%w", err)
	}
	return nil
}

func (c *MyTelnetClient) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("receive:%w", err)
	}
	return nil
}
