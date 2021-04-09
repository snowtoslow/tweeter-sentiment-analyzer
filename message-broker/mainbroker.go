package main

import (
	"log"
	"net"
	"tweeter-sentiment-analyzer/message-broker/broker"
)

func main() {
	var conn net.Conn
	srv := broker.NewBroker(conn)
	log.Fatal(srv.RunBroker(":8088"))
}
