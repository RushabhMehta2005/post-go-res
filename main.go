package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	memstore "github.com/RushabhMehta2005/post-go-res/memstore"
	"github.com/RushabhMehta2005/post-go-res/wal"
)

// The default port on which our application will run
const port = 4242

// Number of concurrent clients which are connected right now
// TODO: Make this variable thread-safe
var concurrent_clients uint8 = 0

// In memory map to store the actual key-value pair data
var store = memstore.NewHashMap()

// Instance of WAL
const WAL_FILE_PATH = "./wal_files/wal_file"
var walWriter, err = wal.NewWAL(WAL_FILE_PATH)

// TODO: Write client library code to interact with our server as well
// TODO: Make data persistent by means of Write-Ahead Logging (WAL)
// TODO: Add checkpointing to WAL
// TODO: Implement new improved open-addressed hash table

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

// TODO: Refactor this function

func handleConnection(conn net.Conn) {
	// No matter how we end up handling this connection, always close the connection and decrease number of clients
	defer func() {
		conn.Close()
		// Possible race on concurrent_clients
		concurrent_clients -= 1
		fmt.Println("Connection closed by client: ", conn.RemoteAddr().String())
		fmt.Println("Concurrent clients: ", concurrent_clients)
	}()

	fmt.Println("Connection established by client: ", conn.RemoteAddr().String())
	fmt.Println("Concurrent clients: ", concurrent_clients)

	reader := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for {
		next := reader.Scan()
		if next {
			cmd_raw_string := reader.Text()
			cmd_fields := strings.Fields(cmd_raw_string)

			// Skip empty newlines by client, QOL change
			if len(cmd_fields) == 0 {
				continue
			}

			switch cmd_fields[0] {
			case "SET":
				switch len(cmd_fields) {
				case 1:
					writer.WriteString("SET NOT OK: No key or value provided\n")
				case 2:
					writer.WriteString("SET NOT OK: No value provided\n")
				case 3:
					walWriter.Log(cmd_fields[0], cmd_fields[1], cmd_fields[2])
					store.Set(cmd_fields[1], cmd_fields[2])
					writer.WriteString("SET OK: Wrote value of " + cmd_fields[1] + " as " + cmd_fields[2] + "\n")
				default:
					writer.WriteString("SET NOT OK: Invalid arguments\n")
				}
			case "GET":
				switch len(cmd_fields) {
				case 1:
					writer.WriteString("GET NOT OK: No key provided\n")
				case 2:
					val, found := store.Get(cmd_fields[1])
					if found {
						writer.WriteString("GET OK: Got value of " + cmd_fields[1] + " as " + val + "\n")
					} else {
						writer.WriteString("GET NOT OK: Did not find " + cmd_fields[1] + " in store" + "\n")
					}
				default:
					writer.WriteString("GET NOT OK: Invalid arguments\n")
				}
			default:
				writer.WriteString("INVALID COMMAND" + "\n")
			}
			err := writer.Flush()
			if err != nil {
				fmt.Println("Could not flush data to client: ", conn.RemoteAddr())
			}

		} else {
			break
		}
	}

}
