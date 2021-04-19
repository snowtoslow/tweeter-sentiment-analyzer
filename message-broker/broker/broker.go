package broker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"tweeter-sentiment-analyzer/message-broker/client"
	"tweeter-sentiment-analyzer/message-broker/parser"
)

type Broker struct {
	connection      net.Conn
	clients         map[string]*client.Client
	counter         int
	ActorConnection net.Conn
	magicChan       chan string //chan of data from client actor
	ClientTopic     map[string][]string
}

func NewBroker(connection net.Conn) *Broker {
	broker := &Broker{
		connection:  connection,
		clients:     make(map[string]*client.Client),
		magicChan:   make(chan string),
		ClientTopic: map[string][]string{},
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

		/*setClientTopics := func(myClient *client.Client,topics parser.Topics) {
			if myClient.SubscribedTopics == nil {
				myClient.SubscribedTopics = &topics
			} else {
				missingElements := missing(*myClient.SubscribedTopics,topics)
				for _, v := range missingElements {
					*(myClient.SubscribedTopics) = append(*(myClient.SubscribedTopics), v)
				}
			}
		}*/

		//refactor here line 82!!!!!!!
		/*deleteClientTopicOnUnsubscriber := func(myClient *client.Client,topics parser.Topics) {
			log.Println("delete")
			for _, topicToUnsubscribe := range topics {
				for positionOfTopicToUnsubscribe, presentTopic := range *myClient.SubscribedTopics {
					if topicToUnsubscribe == presentTopic {
						if positionOfTopicToUnsubscribe >= len(*myClient.SubscribedTopics) || positionOfTopicToUnsubscribe < 0 {
							log.Printf("Index is out of range. Index is %d with slice length %d", positionOfTopicToUnsubscribe, len(*myClient.SubscribedTopics))
							return
						}

						*myClient.SubscribedTopics = append((*myClient.SubscribedTopics)[:positionOfTopicToUnsubscribe], (*myClient.SubscribedTopics)[positionOfTopicToUnsubscribe+1:]...)
						break
					}
				}
				break
			}
		}*/

		writeToConn := func(topic string, client *client.Client) {
			for {
				select {
				case _ = <-broker.magicChan:
					switch command := strings.Split(topic, " ")[0]; command {
					case "subscribe":
						subscribe, err := parseToSubscribe(topic, conn.RemoteAddr().String())
						if err != nil {
							log.Println("Error occurred parsing subscribe:", err)
							return
						}
						//setClientTopics(client,subscribe.Topics)
						broker.ClientTopic[client.Connection.RemoteAddr().String()] = subscribe.Topics
						log.Println("SET:", broker.ClientTopic)
						/*writer := bufio.NewWriter(conn)
						for _, v := range *client.SubscribedTopics {
							if strings.Contains(a,v) {
								//writer := bufio.NewWriter(conn)
								n, err := writer.WriteString(a)
								if err != nil {
									log.Println("write: ", err, n)
									return
								}
								err = writer.Flush()
								if err!=nil {
									log.Println("flush: ", err)
									return
								}
							}
						}*/
					case "unsubscribe":
						var mu sync.Mutex
						unsubscribe, err := parseToUnsubscribeSubscribe(topic, conn.RemoteAddr().String())
						if err != nil {
							log.Println("Error occurred parsing unsubscribe:", err)
							return
						}
						//deleteClientTopicOnUnsubscriber(client,unsubscribe.Topics)
						if v, ok := broker.ClientTopic[conn.RemoteAddr().String()]; ok {
							for _, value := range v {
								for _, value2 := range unsubscribe.Topics {
									if value == value2 {
										mu.Lock()
										log.Println(value2)
										mu.Unlock()
										break
									}
								}
							}
						}
					}

				}
			}
		}

		if v, ok := broker.clients[conn.RemoteAddr().String()]; ok {
			for {
				select {
				case a := <-v.Outgoing:
					log.Println("msg from outgoing:", a)
					go writeToConn(a, v)
				}
			}
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

func missing(a, b parser.Topics) parser.Topics {
	ma := make(map[string]struct{}, len(a))
	var diffs parser.Topics

	for _, ka := range a {
		ma[ka] = struct{}{}
	}

	for _, v := range b {
		if _, ok := ma[v]; !ok {
			diffs = append(diffs, v)
		}
	}

	return diffs
}

func parseToSubscribe(topic, address string) (sub parser.Subscribe, err error) {
	myStr, err := extract(topic)
	if err != nil {
		return
	}
	sub.Address.Addresses = append(sub.Address.Addresses, address)
	err = json.Unmarshal([]byte(myStr), &sub)
	if err != nil {
		return
	}
	return
}

func parseToUnsubscribeSubscribe(topic, address string) (unsub parser.Unsubscribe, err error) {
	myStr, err := extract(topic)
	if err != nil {
		log.Println(err)
		return
	}

	unsub.Address.Addresses = append(unsub.Address.Addresses, address)
	err = json.Unmarshal([]byte(myStr), &unsub)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func extract(topic string) (string, error) {
	regexData := regexp.MustCompile("\\{.*\\:\\{.*\\:.*\\}\\}|\\{(.*?)\\}")
	receivedString := regexData.FindString(topic)
	if receivedString == "" {
		return "", fmt.Errorf("something goes wrong extracting string by regex")
	}
	return receivedString, nil
}

func removeElement(s *parser.Topics, i int) error {
	if i >= len(*s) || i < 0 {
		return fmt.Errorf("Index is out of range. Index is %d with slice length %d", i, len(*s))
	}

	*s = append((*s)[:i], (*s)[i+1:]...)
	return nil
}

func appendCategory(a []string, b []string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}
