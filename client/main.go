package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("server address is mandatory")
	}
	c := NewTCPChatClient(os.Args[1])
	go c.Listen()
	c.Start()
	c.Close()
}
