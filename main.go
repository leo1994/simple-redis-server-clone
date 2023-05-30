package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
		decode, err := parseRedisProtocol(bufio.NewReader(conn))
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			break
		}
		// get first command
		command := decode.Value.([]RedisValue)[0]
		switch command.Value.(string) {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			conn.Write([]byte(fmt.Sprintf("+%s\r\n", decode.Value.([]RedisValue)[1].Value.(string))))
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
