package broker

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"tweeter-sentiment-analyzer/message-broker/client"
)

type Broker struct {
	connection      net.Conn
	clients         map[string]*client.Client
	counter         int
	ActorConnection net.Conn
	magicChan       chan string //chan of data from client actor
}

func NewBroker(connection net.Conn) *Broker {
	broker := &Broker{
		connection: connection,
		clients:    make(map[string]*client.Client),
		magicChan:  make(chan string),
	}
	return broker
}

func (broker *Broker) RunBroker(port string) error {
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
		go broker.handleClients(conn)
	}
}

func (broker *Broker) handleClients(conn net.Conn) {
	broker.counter++
	if broker.counter == 1 {
		broker.ActorConnection = conn
		go broker.readFromActorConnection()
	} else if broker.counter > 1 {
		log.Printf("Client with name: %s connects to broker with remote address of connection: %s", fmt.Sprintf("client_%d", broker.counter), conn.RemoteAddr().String())
		broker.clients[conn.RemoteAddr().String()] = client.NewClient(conn, fmt.Sprintf("client_%d", broker.counter))
		for _, v := range broker.clients {
			v.Listen(broker.magicChan)
		}
	}
}

func (broker *Broker) readFromActorConnection() {
	reader := bufio.NewReader(broker.ActorConnection)
	for {
		line, err := reader.ReadString(10)
		if err != nil {
			log.Println("Error reading from actor:", err)
		}
		broker.magicChan <- line
	}
}
