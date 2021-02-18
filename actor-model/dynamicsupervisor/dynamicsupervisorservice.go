package dynamicsupervisor

import (
	"fmt"
	"log"
	"strconv"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	message_types "tweeter-sentiment-analyzer/actor-model/message-types"
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
		actorPoll = append(actorPoll, workeractor.NewActor("working"+constants.ActorName+strconv.Itoa(i), dynamicSupervisor))
	}
	(*actorregistry.MyActorRegistry)["actorPool"] = &actorPoll

	return
}

func (dynamicSupervisor *DynamicSupervisor) ActorLoop() {
	defer close(dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)
	for {
		select {
		case <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate:
			log.Println("NUMBER!")
			actorNumber := <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate
			if actorNumber == 0 {
				log.Println("SKIP")
				continue
			} else if actorNumber < 0 {
				log.Println("delete", actorNumber)
				dynamicSupervisor.deleteActors(actorNumber)
			} else {
				dynamicSupervisor.addActors(actorNumber)
			}
		case action := <-dynamicSupervisor.ActorProps.ChanToReceiveData:
			dynamicSupervisor.deleteActorByIdentity(action.(message_types.ErrorToSupervisor).ActorIdentity)
			dynamicSupervisor.recreateWorkingActor(action.(message_types.ErrorToSupervisor).ActorIdentity)
		}
	}
}

func (dynamicSupervisor *DynamicSupervisor) SendMessage(msg interface{}) {
	dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate <- msg.(int)
}

func (dynamicSupervisor *DynamicSupervisor) SendErrMessage(msg interface{}) {
	dynamicSupervisor.ActorProps.ChanToReceiveData <- msg
}

func (dynamicSupervisor *DynamicSupervisor) addActors(numberOfActors int) {
	for i := 0; i < numberOfActors; i++ {
		*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
			append(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor),
				workeractor.NewActor("working"+constants.ActorName+strconv.Itoa(i+5), dynamicSupervisor))
	}
}

func (dynamicSupervisor *DynamicSupervisor) deleteActors(numberOfActors int) {
	for i := 0; i < -numberOfActors; i++ {
		*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
			append((*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[:i],
				(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i+1:]...)
	}
}

func (dynamicSupervisor *DynamicSupervisor) deleteActorByIdentity(actorIdentity string) {
	log.Println("Delete actor by identity:", actorIdentity)
	for i := 0; i < len(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor)); i++ {
		if (*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.Identity == actorIdentity {
			*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
				append((*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[:i],
					(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i+1:]...)
		}
	}
}

func (dynamicSupervisor *DynamicSupervisor) recreateWorkingActor(actorName string) {
	log.Println("recreate working actor with identity:", actorName)
	*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
		append(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor),
			workeractor.NewActor(actorName, dynamicSupervisor))
}
