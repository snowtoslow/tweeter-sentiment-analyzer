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

func NewDynamicSupervisor(actorName string) (actorabstraction.IActor, error) {
	chanToReceiveAmountOfActorsToCreate := make(chan int, constants.GlobalChanSize)
	chanRoReceiveErrors := make(chan interface{}, constants.GlobalChanSize)

	dynamicSupervisor := &DynamicSupervisor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanRoReceiveErrors,
		},
		ChanToReceiveNumberOfActorsToCreate: chanToReceiveAmountOfActorsToCreate,
	}

	if err := dynamicSupervisor.CreateActorPoll(constants.DefaultActorPollSize); err != nil {
		return nil, err
	}

	(*actorregistry.MyActorRegistry)["dynamicSupervisor"] = dynamicSupervisor

	go dynamicSupervisor.ActorLoop()

	return dynamicSupervisor, nil
}

func (dynamicSupervisor *DynamicSupervisor) CreateActorPoll(numberOfActors int) (err error) {
	var actorPoll []actorabstraction.IActor
	if numberOfActors <= 1 {
		return fmt.Errorf("number of actors could not be smaller or equal with one")
	}
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, workeractor.NewActor("working_"+strconv.Itoa(i)))
	}
	(*actorregistry.MyActorRegistry)["actorPool"] = &actorPoll

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
				dynamicSupervisor.deleteActors(actorNumber)
			} else {
				dynamicSupervisor.addActors(actorNumber)
			}
		case <-dynamicSupervisor.ActorProps.ChanToReceiveData:
			log.Println("ERROR")
		}
	}
}

func (dynamicSupervisor *DynamicSupervisor) SendMessage(msg interface{}) {
	dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate <- msg.(int)
}

func (dynamicSupervisor *DynamicSupervisor) addActors(numberOfActors int) {
	log.Println("Add actors", numberOfActors)
	for i := 0; i < numberOfActors; i++ {
		*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
			append(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor),
				workeractor.NewActor("working_"+strconv.Itoa(i+5)))
	}
}

func (dynamicSupervisor *DynamicSupervisor) deleteActors(numberOfActors int) {
	log.Println("Delete actors", numberOfActors)
	for i := 0; i < -numberOfActors; i++ {
		*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
			append((*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[:i],
				(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i+1:]...)
	}
}
