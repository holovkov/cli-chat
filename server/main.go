package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("server address is mandatory")
	}
	c := NewTCPChatServer(os.Args[1])
	go c.BroadcastLoop()
	c.Listen()
}
