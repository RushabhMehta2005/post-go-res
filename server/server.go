package server

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/RushabhMehta2005/post-go-res/memstore"
	"github.com/RushabhMehta2005/post-go-res/utils"
	"github.com/RushabhMehta2005/post-go-res/wal"
)

type Server struct {
	kvstore store.InMemStore // Actual in memory store
	walHandler *wal.WAL // Write Ahead Logger
	concurrentClients utils.AtomicCounter // The number of clients connected right now
	port string // The default port on which our application will run
}

func NewServer(kvstore store.InMemStore, walHandler *wal.WAL) (*Server) {
	return &Server{
		kvstore: kvstore,
		walHandler: walHandler,
		concurrentClients: utils.AtomicCounter{},
		port: "4242",
	}
}

func (s *Server) Start() {
	log.Println("Starting post-go-res ...")

	// Rebuild the database in memory from the wal
	s.walHandler.ReBuild(s.kvstore)

	// Listen for tcp connections on the port
	listener, err := net.Listen("tcp", ":" + s.port)

	if err != nil {
		log.Fatal("Could not listen on port", s.port)
	}

	// No matter how the server exits, always clean up the resources
	defer listener.Close()

	log.Println("Server started successfully.")
	log.Println("Listening on port", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Could not accept connection from a client.")
		}

		s.concurrentClients.Inc()

		// Launch a goroutine to handle the connection, allowing main loop to accept more connections
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// No matter how we end up handling this connection, always close the connection and decrease number of clients
	defer func() {
		conn.Close()

		s.concurrentClients.Dec()

		log.Println("Connection closed by client:", conn.RemoteAddr().String())
		log.Println("Concurrent clients:", s.concurrentClients.Read())
	}()

	log.Println("Connection established by client:", conn.RemoteAddr().String())
	log.Println("Concurrent clients:", s.concurrentClients.Read())

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
					s.walHandler.Log(cmd_fields[0], cmd_fields[1], cmd_fields[2])
					s.kvstore.Set(cmd_fields[1], cmd_fields[2])
					writer.WriteString("SET OK: Wrote value of " + cmd_fields[1] + " as " + cmd_fields[2] + "\n")
				default:
					writer.WriteString("SET NOT OK: Invalid arguments\n")
				}
			case "GET":
				switch len(cmd_fields) {
				case 1:
					writer.WriteString("GET NOT OK: No key provided\n")
				case 2:
					val, found := s.kvstore.Get(cmd_fields[1])
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
				log.Println("Could not flush data to client:", conn.RemoteAddr())
			}

		} else {
			break
		}
	}

}
