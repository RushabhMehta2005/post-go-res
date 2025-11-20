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

// Server implements a simple TCP-based key-value server.
//
// The server speaks a tiny text protocol with three commands:
//   - SET <key> <value>\n   -> stores the value and replies "+OK\n"
//   - GET <key>\n         -> returns "+OK <value>\n" or an error "-GET could not find <key> in store\n"
//   - DEL <key>\n         -> deletes the key and replies "+OK\n"
//
// The server writes every mutation (SET/DEL) to a write-ahead-log (WAL) before applying
// it to the in-memory store. On startup the WAL is replayed to rebuild the in-memory state.
//
// Server is safe for concurrent clients: each incoming connection is handled in its own goroutine.
type Server struct {
	kvstore           store.InMemStore    // Actual in memory store
	walHandler        wal.WAL             // Write Ahead Logger
	concurrentClients utils.AtomicCounter // The number of clients connected right now
	port              int                 // The port on which our application will run
}

// NewServer constructs a new Server.
//
// Parameters:
//   - kvstore: implementation of store.InMemStore to use for in-memory data (e.g., HashMap or ShardedMap).
//   - walHandler: WAL instance used for durable logging of mutations.
//   - port: TCP port to listen on
func NewServer(kvstore store.InMemStore, walHandler wal.WAL, port int) *Server {
	return &Server{
		kvstore:           kvstore,
		walHandler:        walHandler,
		concurrentClients: utils.AtomicCounter{},
		port:              port,
	}
}

// Start begins listening on the configured TCP port and accepts incoming client connections.
//
// On startup Start will replay the WAL to rebuild the in-memory store, then it will accept
// connections in a loop. Each accepted connection is handled in a separate goroutine by
// handleConnection.
//
// This method blocks until the listener fails (fatal) or the process exits.
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

// handleConnection runs the request loop for a single client connection.
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

		// Dispatch to the handler for the given command. Handlers return the full
		// text response (including trailing newline) to write back to the client.
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

		// Send response to client.
		writer.WriteString(response)
		err := writer.Flush()
		if err != nil {
			log.Println("Could not flush data to client:", conn.RemoteAddr())
		}
	}
	// scanner.Scan loop ends when the client closes the connection or an error occurs.
}

// handleSet processes the SET command arguments.
//
// Expected args: [key, value].
// The mutation is first appended to the WAL, then applied to the in-memory store.
// Returns a single-line response string (with trailing newline) to send to the client.
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

// handleGet processes the GET command arguments.
//
// Expected args: [key].
// Returns "+OK <value>\n" if the key exists, otherwise an error line indicating the key was not found.
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

// handleDel processes the DEL command arguments.
//
// Expected args: [key].
// The deletion is logged to the WAL before being applied to the in-memory store.
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
