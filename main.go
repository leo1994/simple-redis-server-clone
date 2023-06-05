package main

import (
	"fmt"
	"log"
	"net"

	"github.com/leo1994/simpleRedisServerClone/resp"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Printf("Server is running on port 6379\n")

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
		resp.NewDecoder(conn)
		conn.Write([]byte("+Ok\n\r"))
	}
}
