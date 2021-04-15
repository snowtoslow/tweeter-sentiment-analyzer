package principal

import (
	"bufio"
	"log"
	"net"
)

// MyMagicClient is firstly connected client which is our actor who sends tweets and topics to broker server;
type MyMagicClient struct {
	Outgoing   chan string
	reader     *bufio.Reader
	writer     *bufio.Writer
	Connection net.Conn
	name       string
}

func (mc *MyMagicClient) listen() {
	go mc.read()
}

func (mc *MyMagicClient) read() {
	defer func(Connection net.Conn) {
		err := Connection.Close()
		if err != nil {

		}
	}(mc.Connection)
	for {
		if line, err := mc.reader.ReadString(10); err == nil {
			if mc.Connection != nil {
				//we use here a goroutine because our unbuffered chan block, because there is no a client which read messages from unbuffered chan
				//If the channel is unbuffered, the sender blocks until the receiver has received the value -> from doc
				go func() {
					mc.Outgoing <- line
				}()
				/*client.outgoing <- line <= main case without a separate goroutine for blocked chans which is waiting from reading from; */
			} else {
				break
			}
		} else {
			log.Println("Error occurred reading string in client from connection2: ", err)
			return
		}
	}
}
