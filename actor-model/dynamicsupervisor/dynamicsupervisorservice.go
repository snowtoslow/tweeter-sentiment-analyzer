package dynamicsupervisor

import (
	"fmt"
	"log"
	"strconv"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	message_types "tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
	"tweeter-sentiment-analyzer/actor-model/workeractor"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
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
			actorNumber := <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate
			if actorNumber == 0 {
				continue
			} else if actorNumber < 0 {
				dynamicSupervisor.deleteActors(actorNumber)
			} else {
				dynamicSupervisor.addActors(actorNumber)
			}
		case action := <-dynamicSupervisor.ActorProps.ChanToReceiveData:
			//dynamicSupervisor.deleteActorByIdentity(action.(message_types.ErrorToSupervisor).ActorIdentity)
			/*dynamicSupervisor.recreateWorkingActor(action.(message_types.ErrorToSupervisor).ActorIdentity)*/
			dynamicSupervisor.deleteActorAndRecreateByIdentity(action.(message_types.ErrorToSupervisor).ActorIdentity)
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
	dynamicSupervisor.sendToRouter(dynamicSupervisor.collectDataOfActorsToBeDroppedToSingleChan(-numberOfActors))
	for i := 0; i < -numberOfActors; i++ {
		*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
			append((*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[:i],
				(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i+1:]...)
	}
}

//log.Printf("idetity: %s--->%s",(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.Identity,a.(string))
func (dynamicSupervisor *DynamicSupervisor) deleteActorAndRecreateByIdentity(actorIdentity string) {
	log.Println("Delete actor by identity:", actorIdentity)
	for i := 0; i < len(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor)); i++ {
		if (*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.Identity == actorIdentity {
			recreatedActor := workeractor.NewActor(actorIdentity, dynamicSupervisor)
			a := <-(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.ChanToReceiveData
			recreatedActor.(*workeractor.Actor).SendMessage(a)
			*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
				append((*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[:i],
					(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i+1:]...)
			dynamicSupervisor.pushRecreatedWorkingActorToArray(recreatedActor.(*workeractor.Actor))
			break //maybe necessary to delete!
		}
	}
}

func (dynamicSupervisor *DynamicSupervisor) pushRecreatedWorkingActorToArray(recreatedActor *workeractor.Actor) {
	log.Println("recreate actor with identity:", recreatedActor.ActorProps.Identity)
	*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor) =
		append(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor),
			recreatedActor)
}

func (dynamicSupervisor *DynamicSupervisor) sendToRouter(myChan chan interface{}) {
	go func() {
		for msg := range myChan {
			actorregistry.MyActorRegistry.FindActorByName("routerActor").(*routeractor.RouterActor).SendMessage(msg)
		}
	}()
}

func (dynamicSupervisor *DynamicSupervisor) collectDataOfActorsToBeDroppedToSingleChan(numberOfActors int) chan interface{} {
	var actorsChanToDelete []chan interface{}
	for i := 0; i < numberOfActors; i++ {
		actorsChanToDelete = append(actorsChanToDelete,
			(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.ChanToReceiveData)
	}
	return utils.MergeWaitInterface(actorsChanToDelete)
}
