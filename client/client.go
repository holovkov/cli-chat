package main

import (
	"io"
	"log"
	"net"
	"os"
)

type ChatClient interface {
	Listen()
	Start()
	Close() error
}

// TCPChatClient is implementation of ChatClient that uses TCP as transport
type TCPChatClient struct {
	conn net.Conn
	done chan struct{}
}

func NewTCPChatClient(addr string) *TCPChatClient {
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	return &TCPChatClient{
		conn: conn,
	}
}

func (c *TCPChatClient) Listen() {
	go func() {
		if _, err := io.Copy(os.Stdout, c.conn); err != nil {
			log.Fatal("error reading from server connection")
		}
		c.done <- struct{}{}
	}()
}

func (c *TCPChatClient) Start() {
	if _, err := io.Copy(c.conn, os.Stdin); err != nil {
		c.conn.Close()
		log.Fatal("error sending data to server, shutting down\n")
	}
}

func (c *TCPChatClient) Close() error {
	return c.conn.Close()
}
