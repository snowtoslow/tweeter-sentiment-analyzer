package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	connection  net.Conn
	clients     map[string]*Client
	counter     int
	actorClient *Client
}

func NewServer(connection net.Conn) *Server {
	return &Server{
		connection: connection,
		clients:    make(map[string]*Client),
	}
}

func (server *Server) RunServer(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go server.handleRequest(conn)
	}
}

func (server *Server) handleRequest(conn net.Conn) {
	log.Println("CONECTED!")
	server.counter++
	if server.counter == 1 {
		server.actorClient = NewClient(conn, "actorClient")
		go server.actorClient.read()
	} else if server.counter > 1 {
		log.Println("remote address of new connection:", conn.RemoteAddr().String())
		server.clients[conn.RemoteAddr().String()] = NewClient(conn, fmt.Sprintf("client_%d", server.counter))
		for _, v := range server.clients {
			v.Listen(server.actorClient.outgoing)
		}
	}

}
