package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleClientConnection(conn)
	}
}

func handleClientConnection(conn net.Conn) {
	defer conn.Close()
	for {
		conn.Write([]byte("+PONG\r\n"))
	}
}
