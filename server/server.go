package main

import (
	"log"
	"net"
	"sync"
)

type ChatServer interface {
	Listen()
	BroadcastLoop()
}

// TCPChatServer is implementation of ChatServer interface that uses TCP as transport
type TCPChatServer struct {
	listener net.Listener
	clients  map[client]bool
	messages chan string
	users    map[string]int
	m        sync.RWMutex
}

func NewTCPChatServer(addr string) *TCPChatServer {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return &TCPChatServer{
		listener: listener,
		clients:  make(map[client]bool),
		messages: make(chan string, 1),
		users:    make(map[string]int),
	}
}

// BroadcastLoop launches loop that sends new messages to all of the
// active clients
func (c *TCPChatServer) BroadcastLoop() {
	for {
		select {
		case msg := <-c.messages:
			for cli := range c.clients {
				cli.out <- msg
			}
		default:
		}
	}
}

// Listen initializes loop that accepts and handles new connections
func (c *TCPChatServer) Listen() {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go c.handleConn(conn)
	}
}

// handleConn handles single client connection
func (s *TCPChatServer) handleConn(conn net.Conn) {
	client := s.newClient(conn)
	defer s.closeClient(client)

	for msg := range client.in {
		s.broadcast(client.name + ": " + msg)
	}
}

// newClient initializes client and notify about new user if necessary
func (s *TCPChatServer) newClient(conn net.Conn) client {
	client := initClient(conn)
	client.promptUsername()
	s.addClientForUser(client.name)
	s.notifyIfNewUser(client.name)
	s.clients[client] = true
	return client
}

func (s *TCPChatServer) addClientForUser(name string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.users[name]++
}

func (s *TCPChatServer) closeClient(c client) error {
	s.deleteClientFromUser(c.name)
	s.notifyIfUserLeft(c.name)
	delete(s.clients, c)
	close(c.out)
	return c.conn.Close()
}

func (s *TCPChatServer) deleteClientFromUser(name string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.users[name]--
}

func (s *TCPChatServer) notifyIfNewUser(username string) {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.users[username] == 1 {
		s.broadcast(username + " has arrived")
	}
}

func (s *TCPChatServer) notifyIfUserLeft(username string) {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.users[username] == 0 {
		s.broadcast(username + " has left")
	}
}

func (s *TCPChatServer) broadcast(msg string) {
	s.messages <- msg
}
