package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"tweeter-sentiment-analyzer/message-broker/typedmsg"
	"tweeter-sentiment-analyzer/message-broker/utils"
)

type Client struct {
	CommandChan                            chan typedmsg.Message
	reader                                 *bufio.Reader
	Connection                             net.Conn
	Name                                   string
	SubscribedTopics                       *[]typedmsg.Topic
	UniqueId                               string
	ChanForRemainingMessagesOfDurableTopic chan string
	DurableQueueChan                       chan *typedmsg.DurableQueue
}

func NewClient(connection net.Conn, name string) *Client {
	reader := bufio.NewReader(connection)
	client := &Client{
		CommandChan:                            make(chan typedmsg.Message),
		Connection:                             connection,
		reader:                                 reader,
		Name:                                   name,
		ChanForRemainingMessagesOfDurableTopic: make(chan string),
		DurableQueueChan:                       make(chan *typedmsg.DurableQueue, 1),
	}
	return client
}

func (client *Client) Listen(ch chan string, stopChan chan typedmsg.StopMessage, notifyDurable chan typedmsg.UniqueIdAndAddress) {
	go client.Read(notifyDurable)
	go client.write(ch, stopChan)
}

func (client *Client) Read(notifDurable chan typedmsg.UniqueIdAndAddress) {
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

				if len(messageStruct.UniqueIDForDurable) != 0 {
					info := typedmsg.UniqueIdAndAddress{
						UniqueId:      messageStruct.UniqueIDForDurable,
						ClientAddress: messageStruct.Address,
					}

					go func(info typedmsg.UniqueIdAndAddress) {
						notifDurable <- info
					}(info)
				}

				go func(messageStruct typedmsg.Message) {
					client.CommandChan <- messageStruct
				}(messageStruct)

			} else {
				break
			}
		} else {
			/*log.Println("Error occurred reading string in client from connection: ", err)*/
			return
		}
	}
}

func (client *Client) write(ch chan string, stopChan chan typedmsg.StopMessage) {
	for {
		select {
		case messageStruct := <-client.CommandChan:
			if anonFunc, ok := client.createMapsWithFunction(ch, stopChan)[messageStruct.Command]; ok {
				go anonFunc(messageStruct.Topics)
			}
		case result := <-client.DurableQueueChan:
			topics := utils.ConvertToTopic(result.DurableTopics)
			client.SetSubscribedTopics(topics)
			writer := bufio.NewWriter(client.Connection)
			for _, val := range result.Queue {
				writer.WriteString(val)
				writer.Flush()
			}
		}
	}
}

func (client *Client) createMapsWithFunction(ch chan string, stopChan chan typedmsg.StopMessage) map[typedmsg.Command]func(topics []typedmsg.Topic) {
	mapsWithFunction := map[typedmsg.Command]func(topics []typedmsg.Topic){
		"subscribe": func(topics []typedmsg.Topic) {
			log.Println("subscribe")
			client.SetSubscribedTopics(topics)
			go func() {
				if err := client.writeToClient(ch); err != nil {
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

func (client *Client) createStopMessage() (stopMsg typedmsg.StopMessage) {
	durableTopics, err := client.extractOnlyDurableTopics()
	if err == nil {
		stopMsg.UniqueClientId = generateUuidgen()
		stopMsg.OnlyDurableTopics = durableTopics
		stopMsg.ClientAddress = client.Connection.RemoteAddr().String()
		stopMsg.MyMagicChan = client.ChanForRemainingMessagesOfDurableTopic
		stopMsg.Name = client.Name
	} else {
		log.Println("Extract only durable topics error: ", err)
		stopMsg.OnlyDurableTopics = nil
	}
	return
}

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

func (client *Client) SetSubscribedTopics(topics []typedmsg.Topic) (newSubscribedTopics *[]typedmsg.Topic) {
	if client.SubscribedTopics == nil {
		client.SubscribedTopics = &topics
	} else {
		missingElements := missing(*client.SubscribedTopics, topics)
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
	client.sendRemainingMessagesOfADurableTopic(ch)
	for {
		select {
		case msg := <-client.ChanForRemainingMessagesOfDurableTopic:
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

func (client *Client) sendRemainingMessagesOfADurableTopic(ch chan string) {
	go func() {
		for msg := range ch {
			for _, v := range *client.SubscribedTopics {
				if strings.Contains(msg, v.Value) {
					client.ChanForRemainingMessagesOfDurableTopic <- msg
				}
			}
		}
		close(client.ChanForRemainingMessagesOfDurableTopic)
	}()
}

func (client *Client) deleteTopic(topicsToUnsubscribe []typedmsg.Topic) {
	for _, topicToUnsubscribe := range topicsToUnsubscribe {
		for _, value := range *client.SubscribedTopics {
			if value == topicToUnsubscribe {
				log.Println("function to delete topic here")
			}
		}
	}
}

func generateUuidgen() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatal("ERROR OPENING FILE " + "/dev/urandom")
	}
	b := make([]byte, 16)
	if _, err = f.Read(b); err != nil {
		log.Fatal("ERROR READING FROM FILE: " + "/dev/urandom")
	}
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func missing(a, b []typedmsg.Topic) (diffs []typedmsg.Topic) {
	// create map with length of the 'a' slice
	ma := make(map[string]struct{}, len(a))

	// Convert first slice to map with empty struct (0 bytes)
	for _, ka := range a {
		ma[ka.Value] = struct{}{}
	}
	// find missing values in a
	for _, kb := range b {
		if _, ok := ma[kb.Value]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}
