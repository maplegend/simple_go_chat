package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	fmt.Println("Error!")
	fmt.Println(err)
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
		}
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	reader := bufio.NewReader(client)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
			return
		}
		msgs <- Message{clientid, msg}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		panic(fmt.Sprintf("Cant listen on port %s\n%v", *portPtr, err))
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	lastId := 0

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			clientId := lastId
			lastId ++
			clients[clientId] = conn
			go handleClient(conn, clientId, msgs)

			fmt.Println("New client connected, id:", clientId)
		case msg := <-msgs:
			for id, conn := range clients {
				if msg.sender != id {
					_, err := conn.Write([]byte(msg.message))
					if err != nil {
						handleError(err)
					}
				}
			}
		}
	}
}
