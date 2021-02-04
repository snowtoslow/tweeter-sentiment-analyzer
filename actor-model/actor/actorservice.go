package actor

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	msgType "tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
)

func CreateActorPoll(numberOfActors int) (actorPoll []*Actor, err error) {
	if numberOfActors <= 1 {
		return nil, fmt.Errorf("number of actors could not be smaller or equal with one!")
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
		if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct {
			log.Println("Stuff to count:")
		} else if fmt.Sprintf("%T", action) == constants.PanicMessageType {
			log.Println("ERROR:")
		}
	}
}

func (actor *Actor) SendMessage(data string) {
	actor.ChanToReceiveData <- data
}

func (actor *Actor) processReceivedMessage(dataToRecv string) interface{} {
	regexData := regexp.MustCompile(constants.JsonRegex)
	receivedString := regexData.FindString(dataToRecv)
	var tweetMsg *models.MyJsonName
	if receivedString == constants.PanicMessage {
		return msgType.PanicMessage(receivedString)
	} else {
		err := json.Unmarshal([]byte(receivedString), &tweetMsg)
		if err != nil {
			return err
		}
		return tweetMsg
	}
}
