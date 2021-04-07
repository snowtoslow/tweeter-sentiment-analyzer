package main

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	outgoing      chan string
	reader        *bufio.Reader
	writer        *bufio.Writer
	remoteAddress string
	connection    net.Conn
}

func NewClient(conn net.Conn) *Client {
	writerInitializer := bufio.NewWriter(conn)
	readerInitializer := bufio.NewReader(conn)
	return &Client{
		outgoing:      make(chan string, 10),
		reader:        readerInitializer,
		writer:        writerInitializer,
		remoteAddress: conn.RemoteAddr().String(),
		connection:    conn,
	}
}

func (client *Client) read() {
	for {
		log.Println("THERE:", <-client.outgoing)
	}
}
