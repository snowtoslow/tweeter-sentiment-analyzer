package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"tweeter-sentiment-analyzer/constants"
)

type Client struct {
	outgoing   chan string
	reader     *bufio.Reader
	writer     *bufio.Writer
	connection net.Conn
	name       string
}

func NewClient(connection net.Conn, name string) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		outgoing:   make(chan string),
		connection: connection,
		reader:     reader,
		writer:     writer,
		name:       name,
	}

	return client
}

func (client *Client) Listen(ch chan string) {
	go client.read()
	go client.write(ch)
}

func (client *Client) read() {
	log.Println(client.name)
	defer close(client.outgoing)
	for {
		line, err := client.reader.ReadString(10)
		//log.Printf("client:%s -> line last: %v -> len bytest(%d) -> (%s)",client.connection.RemoteAddr(),[]byte(line)[len([]byte(line))-1],len([]byte(line)),line)
		if err == nil {
			if client.connection != nil {
				//we use here a goroutine because our unbuffered chan block, because there is no a client which read messages from unbuffered chan
				//If the channel is unbuffered, the sender blocks until the receiver has received the value -> from doc
				go func() {
					client.outgoing <- line
				}()
				/*client.outgoing <- line <= main case without a separate goroutine for blocked chans which is waiting from reading from; */
			} else {
				break
			}
		}
	}
	err := client.connection.Close()
	if err != nil {
		log.Println("err:", err)
	}
}

func (client *Client) write(ch <-chan string) {
	for data := range client.outgoing {
		if strings.TrimSpace(data) == constants.TweetsTopic || strings.TrimSpace(data) == constants.UserTopic {
			for val := range ch {
				if strings.Contains(val, strings.TrimSpace(data)) {
					client.writer.WriteString(val)
					client.writer.Flush()
				}
			}
		}
	}
}
