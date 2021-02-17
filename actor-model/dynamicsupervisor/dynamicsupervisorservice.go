package dynamicsupervisor

import (
	"fmt"
	"log"
	"strconv"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/workeractor"
	"tweeter-sentiment-analyzer/constants"
)

func NewDynamicSupervisor(actorName string) actorabstraction.IActor {
	chanToReceiveAmountOfActorsToCreate := make(chan int, constants.GlobalChanSize)
	chanRoReceiveErrors := make(chan interface{}, constants.GlobalChanSize)

	dynamicSupervisor := &DynamicSupervisor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanRoReceiveErrors,
		},
		ChanToReceiveNumberOfActorsToCreate: chanToReceiveAmountOfActorsToCreate,
	}

	if err := dynamicSupervisor.CreateActorPoll(5); err != nil {
		log.Println("ERROR IN DYNAMIC SUPERVISOR HERE!", err)
	}

	(*actorregistry.MyActorRegistry)["dynamicSupervisor"] = dynamicSupervisor

	go dynamicSupervisor.ActorLoop()

	return dynamicSupervisor
}

func (dynamicSupervisor *DynamicSupervisor) CreateActorPoll(numberOfActors int) (err error) {
	actorPoll := new([]actorabstraction.IActor)
	if numberOfActors <= 1 {
		return fmt.Errorf("number of actors could not be smaller or equal with one")
	}
	for i := 0; i < numberOfActors; i++ {
		*actorPoll = append(*actorPoll, workeractor.NewActor("working_"+strconv.Itoa(i)))
	}
	(*actorregistry.MyActorRegistry)["actorPool"] = *actorPoll

	return
}

func (dynamicSupervisor *DynamicSupervisor) ActorLoop() {
	defer close(dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)
	for {
		select {
		case <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate:
			actorNumber := <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate
			if actorNumber == 0 {
				log.Println("SKIP")
				continue
			} else if actorNumber < 0 {
				dynamicSupervisor.deleteActors()
			} else {
				dynamicSupervisor.addActors()
			}
		case <-dynamicSupervisor.ActorProps.ChanToReceiveData:
			log.Println("ERROR")
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
