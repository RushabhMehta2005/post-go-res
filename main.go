package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

const port = 4242

func main() {
	// For TCP, UDP and IP networks, if the host is empty or a literal unspecified IP address, as in ":80",
	//  "0.0.0.0:80" or "[::]:80" for TCP and UDP, "", "0.0.0.0" or "::" for IP, the local system is assumed.
	
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))

	if err != nil {
		fmt.Println(err)
		log.Fatal("Could not listen on port", port)
	}

	// No matter how the main function returns, always stop the server and clean up the resources
	defer listener.Close()

	fmt.Println("Server started successfully.")
	fmt.Println("Listening on port", port)
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	fmt.Println("Connection established by client: ", conn.RemoteAddr().String())

	for {
		data := make([]byte, 128)
		_, err := conn.Read(data)
		
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client: ", conn.RemoteAddr().String())
				return
			}
		}
		
		fmt.Println("Received: ", string(data))
		
		conn.Write(data)
		fmt.Println("Echoed data back: ", string(data))
	}
	

}