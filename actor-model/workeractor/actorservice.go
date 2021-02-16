package workeractor

import (
	"fmt"
	"log"
	"regexp"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewActor(actorName string) *Actor {
	chanToRecv := make(chan string, constants.GlobalChanSize)
	actor := &Actor{
		Identity:          actorName + constants.ActorName,
		ChanToReceiveData: chanToRecv,
	}

	go actor.ActorLoop()
	return actor
}

func (actor *Actor) ActorLoop() {
	defer close(actor.ChanToReceiveData)
	for {
		action := actor.processReceivedMessage(<-actor.ChanToReceiveData)
		actionsLog(action)
	}
}

func (actor *Actor) SendMessage(data string) {
	actor.ChanToReceiveData <- data
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

func (actor *Actor) processReceivedMessage(dataToRecv string) (resultDataStructure interface{}) {
	regexData := regexp.MustCompile(constants.JsonRegex)
	if receivedString := regexData.FindString(dataToRecv); len(receivedString) != 0 {
		// be carefully with error from json CreateMessageType
		resultDataStructure = utils.CreateMessageType(receivedString)
	}
	return
}
