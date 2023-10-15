package main

import (
	"bufio"
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
	return &client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.addr, c.timeout)
	return err
}

func (c *client) Close() error {
	return c.conn.Close()
}

// in -> net
func (c *client) Send() error {
	return negotiator(c.in, c.conn)
}

// net -> out
func (c *client) Receive() error {
	return negotiator(c.conn, c.out)
}

func negotiator(r io.Reader, w io.Writer) error {
	scanner := *bufio.NewScanner(r)

	for {
		if !scanner.Scan() {
			return scanner.Err()
		}
		str := scanner.Text()

		_, err := w.Write([]byte(fmt.Sprintf("%s\n", str)))
		if err != nil {
			return err
		}
	}
}
