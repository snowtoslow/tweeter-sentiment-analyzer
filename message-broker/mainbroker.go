package main

import (
	"log"
	"tweeter-sentiment-analyzer/message-broker/broker"
)

func main() {
	//var conn net.Conn
	brokerServer := broker.NewBroker()
	log.Fatal(brokerServer.RunBroker(":8088"))
}
