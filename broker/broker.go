package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	connection  net.Conn
	clients     []*Client
	counter     int
	actorClient *Client
}

func NewServer(connection net.Conn) *Server {
	return &Server{
		connection: connection,
		clients:    []*Client{},
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
		server.actorClient = NewClient(conn, "mainClientSendMsg")
		go server.actorClient.read()
	} else if server.counter > 1 {
		log.Println("here")
		server.clients = append(server.clients, NewClient(conn, fmt.Sprintf("client_%d", len(server.clients))))
		for i := 0; i < len(server.clients); i++ {
			/*go server.clients[i].read()
			go server.clients[i].write(server.actorClient.outgoing)*/
			log.Println(i)
			server.clients[i].Listen(server.actorClient.outgoing)
			//server.clients[i].Listen(server.actorClient.outgoing)
		}
	}

}
