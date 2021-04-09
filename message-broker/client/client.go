package client

import (
	"bufio"
	"net"
	"strings"
	"tweeter-sentiment-analyzer/message-broker/commands"
)

type Client struct {
	Outgoing   chan string
	reader     *bufio.Reader
	writer     *bufio.Writer
	connection net.Conn
	name       string
}

func NewClient(connection net.Conn, name string) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		Outgoing:   make(chan string),
		connection: connection,
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
	defer close(client.Outgoing)
	for {
		line, err := client.reader.ReadString(10)
		//log.Printf("client:%s -> line last: %v -> len bytest(%d) -> (%s)",client.connection.RemoteAddr(),[]byte(line)[len([]byte(line))-1],len([]byte(line)),line)
		if err == nil {
			if client.connection != nil {
				//we use here a goroutine because our unbuffered chan block, because there is no a client which read messages from unbuffered chan
				//If the channel is unbuffered, the sender blocks until the receiver has received the value -> from doc
				go func() {
					client.Outgoing <- line
				}()
				/*client.outgoing <- line <= main case without a separate goroutine for blocked chans which is waiting from reading from; */
			} else {
				break
			}
		}
	}
}

func (client *Client) write(ch <-chan string) {
	for data := range client.Outgoing {
		if strings.TrimSpace(data) == commands.TweetsTopic || strings.TrimSpace(data) == commands.UsersTopic {
			for val := range ch {
				if strings.Contains(val, strings.TrimSpace(data)) {
					client.writer.WriteString(val)
					client.writer.Flush()
				}
			}
		}
	}
}
