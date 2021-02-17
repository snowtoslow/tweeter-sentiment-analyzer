package workeractor

import (
	"fmt"
	"log"
	"regexp"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewActor(actorName string) actorabstraction.IActor {
	chanToRecv := make(chan interface{}, constants.GlobalChanSize)
	actor := &Actor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToRecv,
		},
	}

	go actor.ActorLoop()

	return actor
}

func (actor *Actor) ActorLoop() {
	defer close(actor.ActorProps.ChanToReceiveData)
	for {
		action := actor.processReceivedMessage(<-actor.ActorProps.ChanToReceiveData)
		if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct {
			//log.Println("Stuff to count:")
		} else if fmt.Sprintf("%T", action) == constants.PanicMessageType {
			log.Println("ERROR:", actor.ActorProps.Identity)
		} else {
			//log.Printf("Nil is received!")
		}
		//actionsLog(action)
	}
}

func (actor *Actor) SendMessage(data interface{}) {
	actor.ActorProps.ChanToReceiveData <- data
}

func actionsLog(action interface{}) {
	if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct {
		//log.Println("Stuff to count:")
	} else if fmt.Sprintf("%T", action) == constants.PanicMessageType {
		log.Println("ERROR:")
	} else {
		//log.Printf("Nil is received!")
	}
}

func (actor *Actor) processReceivedMessage(dataToRecv interface{}) (resultDataStructure interface{}) {
	regexData := regexp.MustCompile(constants.JsonRegex)
	if receivedString := regexData.FindString(dataToRecv.(string)); len(receivedString) != 0 {
		// be carefully with error from json CreateMessageType
		resultDataStructure = utils.CreateMessageType(receivedString)
	}
	return
}
