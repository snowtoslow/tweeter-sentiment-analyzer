package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	connection net.Conn
	clients    map[string]*Client
	counter    int
}

func NewServer(connection net.Conn) *Server {
	return &Server{
		connection: connection,
		clients:    map[string]*Client{},
	}
}

func (server *Server) RunServer(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
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
	// Make a buffer to hold incoming data.
	log.Println("CONECTED!")
	server.counter++
	server.clients[conn.RemoteAddr().String()] = NewClient(conn)
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal("Error reading: ", err)
		}
		/*myInput := strings.TrimSpace(netData)

		if strings.TrimSpace(netData) == "topicUsers" {
			log.Println("CAPTURED!")
		}*/
		v, ok := server.clients[conn.RemoteAddr().String()]
		if ok {
			v.outgoing <- netData
			v.read()
		}

		//fmt.Println(conn.RemoteAddr().String())
	}
	conn.Close()
}
