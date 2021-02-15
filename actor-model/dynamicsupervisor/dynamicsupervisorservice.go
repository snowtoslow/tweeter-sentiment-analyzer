package dynamicsupervisor

import (
	"fmt"
	"log"
	"strconv"
	"tweeter-sentiment-analyzer/actor-model/actor"
	"tweeter-sentiment-analyzer/constants"
)

func NewDynamicSupervisor(actorName string) *DynamicSupervisor {
	chanToReceiveAmountOfActorsToCreate := make(chan int, constants.GlobalChanSize)

	dynamicSupervisor := &DynamicSupervisor{
		Identity:                            actorName + constants.ActorName,
		ChanToReceiveNumberOfActorsToCreate: chanToReceiveAmountOfActorsToCreate,
	}

	go dynamicSupervisor.actorLoop()

	return dynamicSupervisor
}

func (dynamicSupervisor *DynamicSupervisor) CreateActorPoll(numberOfActors int) (actorPoll []*actor.Actor, err error) {
	if numberOfActors <= 1 {
		return nil, fmt.Errorf("number of actors could not be smaller or equal with one")
	}
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, actor.NewActor("working_"+strconv.Itoa(i)))
	}
	return
}

func (dynamicSupervisor *DynamicSupervisor) actorLoop() {
	defer close(dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)
	for {
		log.Println("IN DYNAMIC:", <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)
	}
}

func (dynamicSupervisor *DynamicSupervisor) addActors() {

}

func (dynamicSupervisor *DynamicSupervisor) deleteActors() {

}

func (dynamicSupervisor *DynamicSupervisor) createMainActorPool() {}
