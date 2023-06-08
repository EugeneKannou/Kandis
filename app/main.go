package main

import (
	"Kandis/resp"
	"Kandis/storage"
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	fmt.Println("Kandis server started")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	db := storage.NewSafeMap()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleConnection(conn, db)
	}

}

func handleConnection(conn net.Conn, db *storage.SafeMap) {
	defer conn.Close()

	for {
		data, err := resp.DeserializeRESP(bufio.NewReader(conn))
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("Error decoding RESP: ", err.Error())
			return
		}

		command := strings.ToLower(data.Array()[0].String())
		args := data.Array()[1:]
		fmt.Printf("Got %s command\n", command)
		fmt.Printf("%s", data)

		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))

		case "echo":
			var response []string
			for _, arg := range args {
				response = append(response, arg.String())
			}
			conn.Write([]byte(fmt.Sprintf("+%v\r\n", strings.Join(response, " "))))

		case "set":
			if len(args) == 2 {
				key := args[0].String()
				value := args[1]
				db.Write(key, value.Byte())
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("-Expected 2 args, got %v\r\n", len(args))))
			}

		case "get":
			if len(args) == 1 {
				key := args[0].String()
				value := db.Read(key)
				conn.Write([]byte(fmt.Sprintf("+%s\r\n", value)))
			} else {
				conn.Write([]byte(fmt.Sprintf("-Expected 1 arg, got %v\r\n", len(args))))
			}

		default:
			conn.Write([]byte("-unknown command '" + command + "'\r\n"))
		}
	}
}
