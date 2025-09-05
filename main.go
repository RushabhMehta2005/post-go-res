package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// The default port on which our application will run
const port = 4242

// Number of concurrent clients which are connected right now
var concurrent_clients uint8 = 0

func main() {

	// Listen for tcp connections on the port
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatal("Could not listen on port", port)
	}

	// No matter how the main function returns, always stop the server and clean up the resources
	defer listener.Close()

	fmt.Println("Server started successfully.")
	fmt.Println("Listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not accept connection from a client.")
		}

		concurrent_clients += 1
		// Launch a goroutine to handle the connection
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	// No matter how we end up handling this connection, always close the connection and decrease number of clients
	defer func() {
		conn.Close()
		concurrent_clients -= 1
		fmt.Println("Connection closed by client: ", conn.RemoteAddr().String())
		fmt.Println("Concurrent clients: ", concurrent_clients)
	}()

	fmt.Println("Connection established by client: ", conn.RemoteAddr().String())
	fmt.Println("Concurrent clients: ", concurrent_clients)

	for {
		data := make([]byte, 128)
		_, err := conn.Read(data)

		if err != nil {
			if err == io.EOF {
				break
			}
		}

		fmt.Println("Received: ", string(data))
		conn.Write(data)
		fmt.Println("Echoed data back: ", string(data))
	}

}
