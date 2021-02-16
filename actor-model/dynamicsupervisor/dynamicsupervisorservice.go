package dynamicsupervisor

import (
	"fmt"
	"log"
	"strconv"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/workeractor"
	"tweeter-sentiment-analyzer/constants"
)

func NewDynamicSupervisor(actorName string) *DynamicSupervisor {
	chanToReceiveAmountOfActorsToCreate := make(chan int, constants.GlobalChanSize)

	dynamicSupervisor := &DynamicSupervisor{
		Identity:                            actorName + constants.ActorName,
		ChanToReceiveNumberOfActorsToCreate: chanToReceiveAmountOfActorsToCreate,
	}

	(*actorregistry.MyActorRegistry)["dynamicSupervisor"] = dynamicSupervisor

	go dynamicSupervisor.ActorLoop()

	return dynamicSupervisor
}

func (dynamicSupervisor *DynamicSupervisor) CreateActorPoll(numberOfActors int) (err error) {
	actorPoll := new([]workeractor.Actor)
	if numberOfActors <= 1 {
		return fmt.Errorf("number of actors could not be smaller or equal with one")
	}
	for i := 0; i < numberOfActors; i++ {
		*actorPoll = append(*actorPoll, *workeractor.NewActor("working_" + strconv.Itoa(i)))
	}
	(*actorregistry.MyActorRegistry)["actorPool"] = *actorPoll

	return
}

func (dynamicSupervisor *DynamicSupervisor) ActorLoop() {
	defer close(dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)
	for {
		actorNumber := <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate
		if actorNumber == 0 {
			log.Println("SKIP")
			continue
		} else if actorNumber < 0 {
			dynamicSupervisor.deleteActors()
		} else {
			dynamicSupervisor.addActors()
		}
	}
}

func (dynamicSupervisor *DynamicSupervisor) SendMessage(msg interface{}) {
	log.Println("test message!", msg)
}

func (dynamicSupervisor *DynamicSupervisor) addActors() {
	log.Println("Add actors")
}

func (dynamicSupervisor *DynamicSupervisor) deleteActors() {
	log.Println("Delete actors")
}
