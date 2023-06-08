package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		data, err := DeserializeRESP(bufio.NewReader(conn))
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("Error decoding RESP: ", err.Error())
			return
		}

		command := data.Array()[0].String()
		//args := data.Array()[1:]
		fmt.Printf("Got %s command\n", command)

		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			// handle ECHO
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))

			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
