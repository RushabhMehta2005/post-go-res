package server

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/RushabhMehta2005/post-go-res/memstore"
	"github.com/RushabhMehta2005/post-go-res/utils"
	"github.com/RushabhMehta2005/post-go-res/wal"
)

type Server struct {
	kvstore           store.InMemStore    // Actual in memory store
	walHandler        *wal.WAL            // Write Ahead Logger
	concurrentClients utils.AtomicCounter // The number of clients connected right now
	port              int                 // The port on which our application will run
}

func NewServer(kvstore store.InMemStore, walHandler *wal.WAL, port int) *Server {
	return &Server{
		kvstore:           kvstore,
		walHandler:        walHandler,
		concurrentClients: utils.AtomicCounter{},
		port:              port,
	}
}

func (s *Server) Start() {
	log.Println("Starting post-go-res ...")

	// Rebuild the database in memory from the wal
	s.walHandler.ReBuild(s.kvstore)

	// Listen for tcp connections on the port
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))

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

	s.concurrentClients.Inc()

	log.Println("Connection established by client:", conn.RemoteAddr().String())
	log.Println("Concurrent clients:", s.concurrentClients.Read())

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		// Skip empty newlines by client, QOL change
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		var response string

		switch command {
		case "SET":
			response = s.handleSet(args)
		case "GET":
			response = s.handleGet(args)
		case "DEL":
			response = s.handleDel(args)
		default:
			response = "-INVALID COMMAND\n"
		}

		writer.WriteString(response)
		err := writer.Flush()
		if err != nil {
			log.Println("Could not flush data to client:", conn.RemoteAddr())
		}
	}

}

func (s *Server) handleSet(args []string) string {
	if len(args) != 2 {
		return "-SET expected 2 arguments KEY and VALUE\n"
	}
	key, value := args[0], args[1]
	logEntry := wal.NewSetEntry(&key, &value)
	s.walHandler.Log(logEntry)
	s.kvstore.Set(key, value)
	return "+OK\n"
}

func (s *Server) handleGet(args []string) string {
	if len(args) != 1 {
		return "-GET expected 1 argument KEY\n"
	}
	key := args[0]
	value, ok := s.kvstore.Get(key)
	if !ok {
		return "-GET could not find " + key + " in store\n"
	}
	return "+OK " + value + "\n"
}

func (s *Server) handleDel(args []string) string {
	if len(args) != 1 {
		return "-DEL expected 1 argument KEY\n"
	}
	key := args[0]
	logEntry := wal.NewDelEntry(&key)
	s.walHandler.Log(logEntry)
	s.kvstore.Delete(key)
	return "+OK\n"
}
