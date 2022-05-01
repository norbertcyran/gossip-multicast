package main

import (
	"fmt"
	"net"
)

const (
	ConnHost = "localhost"
	ConnPort = "3333"
	ConnType = "tcp"
)

func main() {
	listener, err := net.Listen(ConnType, fmt.Sprintf("%s:%s", ConnHost, ConnPort))
	if err != nil {
		panic(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error closing:", err.Error())
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received request:", string(buf))
	if _, err := conn.Write(buf); err != nil {
		fmt.Println("Error writing:", err.Error())
	}
	if err := conn.Close(); err != nil {
		fmt.Println("Error closing:", err.Error())
	}
}
