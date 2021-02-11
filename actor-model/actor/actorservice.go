package actor

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func CreateActorPoll(numberOfActors int) (actorPoll []*Actor, err error) {
	if numberOfActors <= 1 {
		return nil, fmt.Errorf("number of actors could not be smaller or equal with one")
	}
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, NewActor("working_"+strconv.Itoa(i)))
	}
	return
}

func NewActor(actorName string) *Actor {
	chanToRecv := make(chan string, constants.GlobalChanSize)
	actor := &Actor{
		Identity:          actorName + constants.ActorName,
		ChanToReceiveData: chanToRecv,
	}

	go actor.actorLoop()
	return actor
}

func (actor *Actor) actorLoop() {
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
