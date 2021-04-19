package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"tweeter-sentiment-analyzer/message-broker/typedmsg"
)

type Client struct {
	Outgoing chan string
	reader   *bufio.Reader
	/*writer     *bufio.Writer*/
	Connection       net.Conn
	name             string
	SubscribedTopics *typedmsg.Topics
}

func NewClient(connection net.Conn, name string) *Client {
	/*writer := bufio.NewWriter(connection)*/
	reader := bufio.NewReader(connection)

	client := &Client{
		Outgoing:   make(chan string),
		Connection: connection,
		reader:     reader,
		/*writer:     writer,*/
		name: name,
	}

	return client
}

func (client *Client) Listen(ch chan string) {
	go client.Read()
	go client.write(ch)
}

func (client *Client) Read() {
	defer client.Connection.Close()
	for {
		if line, err := client.reader.ReadString(10); err == nil {
			if client.Connection != nil {
				//we use here a goroutine because our unbuffered chan block, because there is no a client which read messages from unbuffered chan
				//If the channel is unbuffered, the sender blocks until the receiver has received the value -> from doc
				go func() {
					client.Outgoing <- line
				}()
				/*client.outgoing <- line <= main case without a separate goroutine for blocked chans which is waiting from reading from; */
			} else {
				break
			}
		} else {
			log.Println("Error occurred reading string in client from connection: ", err)
			return
		}
	}
}

func (client *Client) write(ch chan string) {
	mapsWithFunction := map[string]func(topic, address string){
		"subscribe": func(topic, address string) {
			log.Println("subscribe")
			subscribe, err := parseToSubscribe(topic, address)
			if err != nil {
				log.Println("Error occurred parsing subscribe:", err)
				return
			}
			client.SetSubscribedTopics(subscribe.Topics)
			out := make(chan string)
			go func() {
				if err = client.writeToClient(ch, out); err != nil {
					log.Println("Error occurred during writing to client => ", err)
					return
				}
			}()

		},
		"unsubscribe": func(topic, address string) {
			log.Println("unsubscribe")
			unsubscribe, err := parseToUnsubscribeSubscribe(topic, client.Connection.RemoteAddr().String())
			if err != nil {
				log.Println("Error occurred parsing unsubscribe:", err)
				return
			}
			go client.deleteTopic(unsubscribe.Topics)
		},
	}
	for {
		select {
		case topic := <-client.Outgoing:
			log.Printf("client with connection address: %s want to %s ", client.Connection.RemoteAddr().String(), strings.Split(strings.TrimSpace(topic), " ")[0])
			command := strings.Split(topic, " ")[0]
			if anonFunc, ok := mapsWithFunction[command]; ok {
				go anonFunc(topic, client.Connection.RemoteAddr().String())
			}
		}
	}
}

func (client *Client) SetSubscribedTopics(topics typedmsg.Topics) *typedmsg.Topics {
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
	myArr := client.SubscribedTopics
	return myArr
}

func (client *Client) writeToClient(ch chan string, out chan string) error {
	clientWriter := bufio.NewWriter(client.Connection)
	//client.writer = clientWriter
	for msg := range ch {
		for _, v := range *client.SubscribedTopics {
			if strings.Contains(msg, v) {
				go func(msg string, v string) {
					out <- msg
				}(msg, v)
				select {
				case a := <-out:
					if _, err := clientWriter.WriteString(a); err != nil {
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
	return nil
}

func (client *Client) deleteTopic(topicsToUnsubscribe typedmsg.Topics) {
	for _, topicToUnsubscribe := range topicsToUnsubscribe {
		for positionOfTopicToUnsubscribe, presentTopic := range *client.SubscribedTopics {
			if topicToUnsubscribe == presentTopic {
				*client.SubscribedTopics = append((*client.SubscribedTopics)[:positionOfTopicToUnsubscribe], (*client.SubscribedTopics)[positionOfTopicToUnsubscribe+1:]...)
			}
		}
	}
}

func missing(a, b typedmsg.Topics) typedmsg.Topics {
	ma := make(map[string]struct{}, len(a))
	var diffs typedmsg.Topics

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

func parseToSubscribe(topic, address string) (sub typedmsg.Subscribe, err error) {
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

func parseToUnsubscribeSubscribe(topic, address string) (unsub typedmsg.Unsubscribe, err error) {
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
	log.Println(topic)
	receivedString := regexData.FindString(topic)
	if receivedString == "" {
		return "", fmt.Errorf("something goes wrong extracting string by regex")
	}
	return receivedString, nil
}
