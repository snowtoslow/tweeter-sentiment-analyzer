package client

import (
	"bufio"
	"log"
	"net"
	"strings"
	"tweeter-sentiment-analyzer/message-broker/commands"
)

type Client struct {
	Outgoing   chan string
	reader     *bufio.Reader
	writer     *bufio.Writer
	Connection net.Conn
	name       string
}

func NewClient(connection net.Conn, name string) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		Outgoing:   make(chan string),
		Connection: connection,
		reader:     reader,
		writer:     writer,
		name:       name,
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
	myBoll := false
	for {
		select {
		case topic := <-client.Outgoing:
			log.Printf("client with connection address: %s want to sbcribe to topic: %s", client.Connection.RemoteAddr().String(), topic)
			myBoll = false
			switch withNoSpace := strings.TrimSpace(topic); withNoSpace {
			case commands.TweetsTopic:
				go func() {
				TWEETSLABEL:
					for msg := range ch {
						if strings.Contains(msg, commands.TweetsTopic) {
							n, err := client.writer.WriteString(msg)
							if err != nil {
								log.Println("write user: ", err, n)
								return
							}
							err = client.writer.Flush()
							if err != nil {
								log.Println("flush user:", err)
								return
							}
						}
						if myBoll {
							break TWEETSLABEL
						}
					}
				}()
			case commands.UsersTopic:
				go func() {
				USERLABEL:
					for msg := range ch {
						if strings.Contains(msg, commands.UsersTopic) {
							n, err := client.writer.WriteString(msg)
							if err != nil {
								log.Println("write user: ", err, n)
								return
							}
							err = client.writer.Flush()
							if err != nil {
								log.Println("flush user:", err)
								return
							}
						}
						if myBoll {
							break USERLABEL
						}
					}
				}()
			case "STOP":
				myBoll = true
				continue
			}
		}
	}
}
