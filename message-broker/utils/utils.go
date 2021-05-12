package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"tweeter-sentiment-analyzer/message-broker/typedmsg"
)

var BrokerClientIdDirectoryName = "/home/snowtoslow/Desktop/broker-client-id"

func CreateFile(dirName, clientID string) error {
	if _, err := os.Stat(BrokerClientIdDirectoryName); !os.IsNotExist(err) {
		//create file!
		err := os.Mkdir(dirName, 0777)
		if err != nil {
			return err
		}

		//create struct client
		clientIdStruct := typedmsg.ClientId{
			Value: clientID,
		}
		jsonBytes, err := json.MarshalIndent(clientIdStruct, "", " ")
		if err != nil {
			log.Fatal("ERROR marshaling ident:", err)
			return err
		}

		//write to created file:
		f, err := os.OpenFile(fmt.Sprintf("%s/%s.json", dirName, clientID), os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			log.Fatal("opening file:", err)
			return err
		}

		_, err = f.Write(jsonBytes)
		if err != nil {
			f.Close()
			log.Fatal("error writing to file!")
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func ConvertToTopic(strings []string) []typedmsg.Topic {
	var topics []typedmsg.Topic
	for _, v := range strings {
		topics = append(topics, typedmsg.Topic{
			Value:     v,
			IsDurable: true,
		})
	}
	return topics
}
