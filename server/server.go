package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func main() {

	listen, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Println("Error while listening", err)
	}
	defer listen.Close()

	for {
		connection, err := listen.Accept()
		if err != nil {
			log.Println("Accept error: ", err)
			return
		}

		go handleConnection(connection)

	}
}

func handleConnection(connection net.Conn) error {
	log.Println("Handle connection is started!")

	for {
		netData, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			log.Fatal("Error reading: ", err)
			return err
		}

		myInput := strings.TrimSpace(netData)

		if myInput == "STOP" {
			break
		}

		// connection.Write([]byte(string(myInput))) -> for writing to client actor in process of find topic
	}

	err := connection.Close()
	if err != nil {
		log.Println("Server close error:", err)
		return err
	}

	return nil
}
