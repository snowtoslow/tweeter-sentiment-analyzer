package broker

import (
	"fmt"
	"log"
	"net"
	"os"
	"tweeter-sentiment-analyzer/message-broker/client"
)

type Broker struct {
	connection  net.Conn
	clients     map[string]*client.Client
	counter     int
	actorClient *client.Client
}

func NewBroker(connection net.Conn) *Broker {
	return &Broker{
		connection: connection,
		clients:    make(map[string]*client.Client),
	}
}

func (server *Broker) RunBroker(port string) error {
	log.Println("Start broker:")
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
		go server.handleClients(conn)
	}
}

func (server *Broker) handleClients(conn net.Conn) {
	server.counter++
	if server.counter == 1 {
		server.actorClient = client.NewClient(conn, "actorClient")
		go server.actorClient.Read()
	} else if server.counter > 1 {
		log.Printf("Client with name: %s connects to broker with remote address of connection: %s", fmt.Sprintf("client_%d", server.counter), conn.RemoteAddr().String())
		server.clients[conn.RemoteAddr().String()] = client.NewClient(conn, fmt.Sprintf("client_%d", server.counter))
		for _, v := range server.clients {
			v.Listen(server.actorClient.Outgoing)
		}
	}

}
