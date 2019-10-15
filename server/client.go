package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// client is single chat connection
type client struct {
	out  chan string
	in   chan string
	name string
	conn net.Conn
}

// InitClient initializes client and launches inbound and outbound listeners
// in new goroutines
func initClient(conn net.Conn) client {
	client := client{
		out:  make(chan string, 1),
		in:   make(chan string),
		conn: conn,
	}
	go client.writer(conn)
	go client.reader(conn)
	return client
}

// promptUsername asks user to input his username and notifies him about it
func (c *client) promptUsername() {
	c.out <- "Enter your name:"
	c.name = <-c.in
	c.out <- "Your name is " + c.name
}

func (c *client) writer(conn net.Conn) {
	for msg := range c.out {
		if _, err := fmt.Fprintln(conn, msg); err != nil {
			// also we could consider closing connection for faulty clients
			log.Printf("error sending message to client addr %s: %s", c.conn.RemoteAddr().String(), err.Error())
		}
	}
}

func (c *client) reader(conn net.Conn) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		c.in <- input.Text()
	}
	if err := input.Err(); err != nil {
		log.Printf("error receiving message from client addr %s: %s", c.conn.RemoteAddr().String(), err.Error())
	}
	close(c.in)
}
