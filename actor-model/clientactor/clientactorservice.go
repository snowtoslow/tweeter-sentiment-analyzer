package clientactor

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/constants"
	"tweeter-sentiment-analyzer/actor-model/models"
)

func NewClientActor(actorName string) actorabstraction.IActor {
	chanForMessages := make(chan interface{}, constants.GlobalChanSize)
	dialer := new(net.Dialer)

	clientActor := &ClientActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanForMessages,
		},
		Connection: *dialer,
	}

	(*actorregistry.MyActorRegistry)["clientActor"] = clientActor

	go clientActor.ActorLoop()

	return clientActor
}

func (clientActor *ClientActor) ActorLoop() {
	defer close(clientActor.ActorProps.ChanToReceiveData)
	conn, err := clientActor.Connection.Dial("tcp", os.Getenv("BROKER_URL")) // change here from "localhost:8088"
	if err != nil {
		log.Println("Error during connection to message-broker: ", err)
		return
	}
	defer conn.Close() // maybe budet kakaeato xueta
	for {
		select {
		case action := <-clientActor.ActorProps.ChanToReceiveData:
			if err = clientActor.sendBrokerMessageToBroker(action, conn); err != nil {
				log.Printf("Error during writing messages to message-broker: %s", err)
				return
			}
		}
	}
}

func (clientActor *ClientActor) sendBrokerMessageToBroker(action interface{}, conn net.Conn) (err error) {
	brokerMsg := new(models.BrokerMsg)
	brokerMsg.Content = action
	brokerMsg.SetTopic(fmt.Sprintf("%T", action))
	out, err := json.Marshal(brokerMsg)
	if err != nil {
		log.Println("ERROR DURING MARSHALLING IN :", err)
		return
	}
	out = append(out, 10)
	if _, err = conn.Write(out); err != nil {
		log.Println("Error during writing to server: ", err)
		return
	} else {
		log.Println("msg send")
	}
	return
}

func (clientActor *ClientActor) SendMessage(msg interface{}) {
	clientActor.ActorProps.ChanToReceiveData <- msg
}
