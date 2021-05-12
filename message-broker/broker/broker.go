package broker

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"tweeter-sentiment-analyzer/message-broker/client"
	"tweeter-sentiment-analyzer/message-broker/typedmsg"
)

type Broker struct {
	/*connection          net.Conn*/
	clients             map[string]*client.Client
	counter             int
	ActorConnection     net.Conn
	actorDataChan       chan string //chan of data from client actor
	durableClientQueues map[string]*typedmsg.DurableQueue
	stopSignalChan      chan typedmsg.StopMessage
	notifyDurable       chan typedmsg.UniqueIdAndAddress
}

func NewBroker(connection net.Conn) *Broker {
	broker := &Broker{
		/*connection:          connection,*/
		clients:             make(map[string]*client.Client),
		actorDataChan:       make(chan string),
		durableClientQueues: make(map[string]*typedmsg.DurableQueue),
		stopSignalChan:      make(chan typedmsg.StopMessage, 1),
		notifyDurable:       make(chan typedmsg.UniqueIdAndAddress, 1),
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
		go func() {
			err = broker.handleClients(conn)
			if err != nil {
				return
			}
		}()
	}
}

func (broker *Broker) handleClients(conn net.Conn) (err error) {
	broker.counter++
	if broker.counter == 1 {
		broker.ActorConnection = conn
		go broker.readFromActorConnection()
	} else if broker.counter > 1 {
		log.Printf("Client with name: %s connects to broker with remote address of connection: %s", fmt.Sprintf("client_%d", broker.counter), conn.RemoteAddr().String())
		broker.clients[conn.RemoteAddr().String()] = client.NewClient(conn, fmt.Sprintf("client_%d", broker.counter))
		for _, v := range broker.clients {
			v.Listen(broker.actorDataChan, broker.stopSignalChan, broker.notifyDurable)
		}

		err = broker.durableMessages()
		if err != nil {
			return err
		}

	}
	return nil
}

func (broker *Broker) durableMessages() error {
	select {
	case stopCommand := <-broker.stopSignalChan:
		log.Println("STOP:", stopCommand)
		v, ok := broker.clients[stopCommand.ClientAddress]
		if !ok {
			log.Printf("client with uniqueId %s not found", stopCommand.ClientAddress)
			return fmt.Errorf("client with uniqueId %s not found", stopCommand.ClientAddress)
		}

		if stopCommand.OnlyDurableTopics != nil && len(stopCommand.OnlyDurableTopics) != 0 {
			queue := typedmsg.NewDurableQueue(stopCommand.UniqueClientId, stopCommand.OnlyDurableTopics)
			for _, topic := range stopCommand.OnlyDurableTopics {
				go func(topic string) {
					for msg := range v.ChanForRemainingMessagesOfDurableTopic {
						if strings.Contains(msg, topic) {
							queue.Enqueue(msg)
						}
					}
				}(topic)
			}
			broker.durableClientQueues[stopCommand.UniqueClientId] = queue
		}

		v.Connection = nil
		delete(broker.clients, stopCommand.ClientAddress)
	case a := <-broker.notifyDurable:
		if queueToSend, ok := broker.durableClientQueues[a.UniqueId]; ok {
			foundClient, err := broker.getClientByAddress(string(a.ClientAddress))
			if err != nil {
				return err
			}
			foundClient.DurableQueueChan <- queueToSend

		} else {
			return fmt.Errorf("client wasn't found!")
		}
	}
	return nil
}

func (broker *Broker) getClientByAddress(address string) (*client.Client, error) {
	if foundClient, ok := broker.clients[address]; ok {
		return foundClient, nil
	}
	return nil, fmt.Errorf("client not exists")
}

func (broker *Broker) readFromActorConnection() {
	reader := bufio.NewReader(broker.ActorConnection)
	for {
		line, err := reader.ReadString(10)
		if err != nil {
			log.Println("Error reading from actor:", err)
		}
		broker.actorDataChan <- line
	}
}
