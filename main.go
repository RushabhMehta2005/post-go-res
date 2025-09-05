package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// The default port on which our application will run
const port = 4242

// Number of concurrent clients which are connected right now
var concurrent_clients uint8 = 0

// In memory map to store the actual key-value pair data
var store map[string]string = make(map[string]string, 32)

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
		// Launch a goroutine to handle the connection, allowing main loop to accept more connections
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

	scanner := bufio.NewScanner(conn)
	for {
		next := scanner.Scan()
		if next {
			cmd_raw_string := scanner.Text()
			cmd_fields := strings.Fields(cmd_raw_string)
			switch cmd_fields[0] {
			case "SET": // length better be 3
				store[cmd_fields[1]] = cmd_fields[2]
				conn.Write([]byte("SET OK: Wrote value of " + cmd_fields[1] + " as " + cmd_fields[2] + "\n"))
			case "GET": // length better be 2
				val, ok := store[cmd_fields[1]]
				if ok {
					conn.Write([]byte("GET OK: Got value of " + cmd_fields[1] + " as " + val + "\n"))
				} else {
					conn.Write([]byte("GET NOT OK: Did not find " + cmd_fields[1] + " in store" + "\n"))
				}
			default:
				conn.Write([]byte("INVALID COMMAND" + "\n"))
			}

			// Only to see
			fmt.Println(store)
		} else {
			break
		}
	}

}
