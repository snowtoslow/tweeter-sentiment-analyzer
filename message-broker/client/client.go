package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"tweeter-sentiment-analyzer/message-broker/typedmsg"
	"tweeter-sentiment-analyzer/message-broker/utils"
)

type Client struct {
	//CommandChan is a chan for receiving commands like subscribe, unsubscribe and stop;
	CommandChan chan typedmsg.Message

	//reader is a created for reading from connection
	reader *bufio.Reader

	//Connection broker connection which contains information about local, remote address and other stuff about tcp conn;
	Connection net.Conn

	//Name is the name of connected client -> client_%nr_of_connected_actor;
	Name string

	//SubscribedTopics topics to which is subscribed a specific client;
	SubscribedTopics *[]typedmsg.Topic

	//UniqueId for client durable topics;
	UniqueId string

	//ClientMessageChanRelatedToATopic chan where messages of subscribed topic is flying;
	ClientMessageChanRelatedToATopic chan string

	//DurableQueueChan is a chan where is is send durable queue, the list of previously subscribed topics;
	DurableQueueChan chan *typedmsg.DurableQueue
}

//NewClient function which creates a new client which connects to broker;
func NewClient(connection net.Conn, name string) *Client {
	reader := bufio.NewReader(connection)
	client := &Client{
		CommandChan:                      make(chan typedmsg.Message),
		Connection:                       connection,
		reader:                           reader,
		Name:                             name,
		ClientMessageChanRelatedToATopic: make(chan string),
		DurableQueueChan:                 make(chan *typedmsg.DurableQueue, 1),
	}
	return client
}

//Listen method which contains parallel read and write operations;
func (client *Client) Listen(ch chan string, stopChan chan typedmsg.StopMessage, notifyDurable chan typedmsg.UniqueIdAndAddress) {
	go client.read(notifyDurable)
	go client.write(ch, stopChan)
}

func (client *Client) read(notifyAboutDurable chan typedmsg.UniqueIdAndAddress) {
	for {
		if line, err := client.reader.ReadString('\n'); err == nil {
			if client.Connection != nil {
				/*we use here a goroutine because our unbuffered chan block, because there is no a client which read messages from unbuffered chan
				If the channel is unbuffered, the sender blocks until the receiver has received the value -> from doc*/
				var messageStruct typedmsg.Message
				if err := json.Unmarshal([]byte(line), &messageStruct); err != nil {
					log.Println("Error occurred unmarshalling string into message struct: ", err)
					return
				}
				messageStruct.Address = typedmsg.ClientAddress(client.Connection.RemoteAddr().String())

				//here is the check if message struct contains uniqueID of a previously subscribed topic
				//
				// and if it exists we notify broker to firstly write messages durable messages;
				if len(messageStruct.UniqueIDForDurable) != 0 {
					notifyBrokerAboutDurableMessages(notifyAboutDurable, messageStruct)
				}

				//default case to receive a command of subscribe or unsubscribe;
				go func(messageStruct typedmsg.Message) {
					client.CommandChan <- messageStruct
				}(messageStruct)

			} else {
				break
			}
		} else {
			return
		}
	}
}

//notifyBrokerAboutDurableMessages function to notify broker that this client was already connected and have durable topics;
func notifyBrokerAboutDurableMessages(notifyAboutDurable chan typedmsg.UniqueIdAndAddress, messageStruct typedmsg.Message) {
	info := typedmsg.UniqueIdAndAddress{
		UniqueId:      messageStruct.UniqueIDForDurable,
		ClientAddress: messageStruct.Address,
	}

	go func(info typedmsg.UniqueIdAndAddress) {
		notifyAboutDurable <- info
	}(info)
}

func (client *Client) write(actorChanWithMessages chan string, stopChan chan typedmsg.StopMessage) {
	for {
		select {
		case messageStruct := <-client.CommandChan:
			if anonFunc, ok := client.createMapsWithFunction(actorChanWithMessages, stopChan)[messageStruct.Command]; ok {
				go anonFunc(messageStruct.Topics)
			}
		case result := <-client.DurableQueueChan:
			topics := utils.ConvertToTopic(result.DurableTopics)
			client.setSubscribedTopics(topics)
			writer := bufio.NewWriter(client.Connection)
			for _, val := range result.Queue {
				writer.WriteString(val)
				writer.Flush()
			}
		}
	}
}

func (client *Client) createMapsWithFunction(actorChanWithMessages chan string, stopChan chan typedmsg.StopMessage) map[typedmsg.Command]func(topics []typedmsg.Topic) {
	mapsWithFunction := map[typedmsg.Command]func(topics []typedmsg.Topic){
		"subscribe": func(topics []typedmsg.Topic) {
			log.Println("subscribe")
			client.setSubscribedTopics(topics)
			go func() {
				if err := client.writeToClient(actorChanWithMessages); err != nil {
					log.Println("Error occurred during writing to client => ", err)
					return
				}
			}()
		},
		"unsubscribe": func(topics []typedmsg.Topic) {
			log.Println("unsubscribe")
			go client.deleteTopic(topics)
		},
		"stop": func(topics []typedmsg.Topic) {
			close(client.CommandChan)
			if err := client.Connection.Close(); err != nil {
				log.Println("Error occurred closing client connection:", err)
			}
			stopChan <- client.createStopMessage()
		},
	}

	return mapsWithFunction
}

//createStopMessage method to create a stop message which will be send to broker chan with stop command;
func (client *Client) createStopMessage() (stopMsg typedmsg.StopMessage) {
	durableTopics, err := client.extractOnlyDurableTopics()
	if err == nil {
		stopMsg.UniqueClientId = utils.GenerateUuid()
		stopMsg.OnlyDurableTopics = durableTopics
		stopMsg.ClientAddress = client.Connection.RemoteAddr().String()
		stopMsg.MyMagicChan = client.ClientMessageChanRelatedToATopic
		stopMsg.Name = client.Name
	} else {
		log.Println("Extract only durable topics error: ", err)
		stopMsg.OnlyDurableTopics = nil
	}
	return
}

//extractOnlyDurableTopics method which extracts only durable topics from topics to which a client is subscribed;
func (client *Client) extractOnlyDurableTopics() (onlyDurableTopics typedmsg.DurableTopicsValue, err error) {
	for _, v := range *client.SubscribedTopics {
		if v.IsDurable {
			onlyDurableTopics = append(onlyDurableTopics, v.Value)
		}
	}

	if len(onlyDurableTopics) == 0 {
		err = fmt.Errorf("client with address %s doesn't have durable topics", client.Connection.RemoteAddr().String())
	}

	return
}

//setSubscribedTopics method which set the subscribed topics of a client;
func (client *Client) setSubscribedTopics(topics []typedmsg.Topic) (newSubscribedTopics *[]typedmsg.Topic) {
	if client.SubscribedTopics == nil {
		client.SubscribedTopics = &topics
	} else {
		missingElements := utils.Missing(*client.SubscribedTopics, topics)
		if missingElements != nil {
			for _, v := range missingElements {
				*(client.SubscribedTopics) = append(*(client.SubscribedTopics), v)
			}
		}
	}
	newSubscribedTopics = client.SubscribedTopics
	return
}

func (client *Client) writeToClient(ch chan string) error {
	clientWriter := bufio.NewWriter(client.Connection)
	client.sendMessagesOfSubscribedTopicsToClientChan(ch)
	for {
		select {
		case msg := <-client.ClientMessageChanRelatedToATopic:
			if client.Connection != nil {
				if _, err := clientWriter.WriteString(msg); err != nil {
					log.Printf("Error writing because of %s ", err)
					return err
				}
				if err := clientWriter.Flush(); err != nil {
					log.Printf("Error flushing because of %s ", err)
					return err
				}
			}
		}
	}
}

//sendMessagesOfSubscribedTopicsToClientChan send messages of subscribed topics to client chan;
func (client *Client) sendMessagesOfSubscribedTopicsToClientChan(ch chan string) {
	go func() {
		for msg := range ch {
			for _, v := range *client.SubscribedTopics {
				if strings.Contains(msg, v.Value) {
					client.ClientMessageChanRelatedToATopic <- msg
				}
			}
		}
		close(client.ClientMessageChanRelatedToATopic)
	}()
}

func (client *Client) deleteTopic(topicsToUnsubscribe []typedmsg.Topic) {
	log.Println("DELETE!")
	for _, topicToUnsubscribe := range topicsToUnsubscribe {
		for key, value := range *client.SubscribedTopics {
			if value.Value == topicToUnsubscribe.Value && !value.IsDurable {
				*client.SubscribedTopics = append((*client.SubscribedTopics)[:key], (*client.SubscribedTopics)[key+1:]...)
			}
		}
	}
}
