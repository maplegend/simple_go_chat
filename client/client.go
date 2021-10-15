package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

func read(conn *net.Conn) {
	reader := bufio.NewReader(*conn)

	fmt.Println("Read worker started")

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading")
			fmt.Println(err)
			return
		}
		fmt.Print(msg)
	}
}

func write(conn *net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Writer worker started")

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Reading error")
			fmt.Println(err)
		}
		_, err = (*conn).Write([]byte(text))
		if err != nil {
			fmt.Println("Send error")
			fmt.Println(err)
		}
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()

	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		fmt.Println("Connection error")
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to server")

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		write(&conn)
		wg.Done()
	}()

	go func() {
		read(&conn)
		wg.Done()
	}()

	fmt.Println("Start workers")

	wg.Wait()

	fmt.Println("One thread exited, exiting client")
}
