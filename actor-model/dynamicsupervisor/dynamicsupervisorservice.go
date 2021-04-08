package dynamicsupervisor

import (
	"fmt"
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

	if err := dynamicSupervisor.CreateActorPoll(constants.DefaultActorPollSize, constants.SentimentActorPool); err != nil {
		return nil, err
	}

	if err := dynamicSupervisor.CreateActorPoll(constants.DefaultActorPollSize, constants.AggregationActorPool); err != nil {
		return nil, err
	}

	(*actorregistry.MyActorRegistry)["dynamicSupervisor"] = dynamicSupervisor

	go dynamicSupervisor.ActorLoop()

	return dynamicSupervisor, nil
}

func (dynamicSupervisor *DynamicSupervisor) CreateActorPoll(numberOfActors int, actorPollName string) (err error) {
	var actorPoll []actorabstraction.IActor
	if numberOfActors <= 1 {
		return fmt.Errorf("number of actors could not be smaller or equal with one")
	}
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, workeractor.NewActor(actorPollName+"/working"+constants.ActorName+strconv.Itoa(i), dynamicSupervisor))
	}
	(*actorregistry.MyActorRegistry)[actorPollName] = &actorPoll
	return
}

func (dynamicSupervisor *DynamicSupervisor) ActorLoop() {
	defer close(dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)
	defer close(dynamicSupervisor.ActorProps.ChanToReceiveData)
	for {
		select {
		case actorNumber := <-dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate:
			if actorNumber == 0 {
				continue
			} else if actorNumber < 0 {
				dynamicSupervisor.deleteActorsFromBothPools(actorNumber)
			} else {
				dynamicSupervisor.addActorsToBothPools(actorNumber)
			}
		case action := <-dynamicSupervisor.ActorProps.ChanToReceiveData:
			go func() {
				dynamicSupervisor.deleteActorAndRecreateByIdentity(action.(message_types.ErrorToSupervisor).ActorIdentity)
			}()
		}
	}
}

func (dynamicSupervisor *DynamicSupervisor) deleteActorsFromBothPools(numberOfActors int) {
	poolsName := []string{constants.AggregationActorPool, constants.SentimentActorPool}
	for _, v := range poolsName {
		go func(v string) {
			dynamicSupervisor.deleteActors(numberOfActors, v)
		}(v)
	}
}

func (dynamicSupervisor *DynamicSupervisor) addActorsToBothPools(numberOfActors int) {
	poolsName := []string{constants.AggregationActorPool, constants.SentimentActorPool}
	for _, v := range poolsName {
		go func(v string) {
			dynamicSupervisor.addActors(numberOfActors, v)
		}(v)
	}
}

func (dynamicSupervisor *DynamicSupervisor) SendMessage(msg interface{}) {
	dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate <- msg.(int)
}

func (dynamicSupervisor *DynamicSupervisor) SendErrMessage(msg interface{}) {
	dynamicSupervisor.ActorProps.ChanToReceiveData <- msg
}

//reuse this function for two arrays and reuse it
func (dynamicSupervisor *DynamicSupervisor) addActors(numberOfActors int, poolName string) {
	for i := 0; i < numberOfActors; i++ {
		*actorregistry.MyActorRegistry.FindActorByName(poolName).(*[]actorabstraction.IActor) =
			append(*actorregistry.MyActorRegistry.FindActorByName(poolName).(*[]actorabstraction.IActor),
				workeractor.NewActor(poolName+"/working"+constants.ActorName+strconv.Itoa(i+5), dynamicSupervisor))
	}
}

//reuse this function for both array and do it concurrently;
func (dynamicSupervisor *DynamicSupervisor) deleteActors(numberOfActors int, poolName string) {
	dynamicSupervisor.sendToRouter(dynamicSupervisor.collectDataOfActorsToBeDroppedToSingleChan(-numberOfActors, poolName))
	for i := 0; i < -numberOfActors; i++ {
		*actorregistry.MyActorRegistry.FindActorByName(poolName).(*[]actorabstraction.IActor) =
			append((*actorregistry.MyActorRegistry.FindActorByName(poolName).(*[]actorabstraction.IActor))[:i],
				(*actorregistry.MyActorRegistry.FindActorByName(poolName).(*[]actorabstraction.IActor))[i+1:]...)
	}
}

//log.Printf("idetity: %s--->%s",(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.Identity,a.(string))
func (dynamicSupervisor *DynamicSupervisor) deleteActorAndRecreateByIdentity(actorIdentity string) {
	//log.Println("Delete actor by identity:", actorIdentity)
	concreteActorPool := utils.GetActorPollNameByActorIdentity(actorIdentity)
	for i := 0; i < len(*actorregistry.MyActorRegistry.FindActorByName(concreteActorPool).(*[]actorabstraction.IActor)); i++ {
		if (*actorregistry.MyActorRegistry.FindActorByName(concreteActorPool).(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.Identity == actorIdentity {
			recreatedActor := workeractor.NewActor(actorIdentity, dynamicSupervisor)

			a := <-(*actorregistry.MyActorRegistry.FindActorByName(concreteActorPool).(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.ChanToReceiveData

			recreatedActor.(*workeractor.Actor).SendMessage(a)

			*actorregistry.MyActorRegistry.FindActorByName(concreteActorPool).(*[]actorabstraction.IActor) =
				append((*actorregistry.MyActorRegistry.FindActorByName(concreteActorPool).(*[]actorabstraction.IActor))[:i],
					(*actorregistry.MyActorRegistry.FindActorByName(concreteActorPool).(*[]actorabstraction.IActor))[i+1:]...)

			dynamicSupervisor.pushRecreatedWorkingActorToArray(recreatedActor.(*workeractor.Actor))

			break //maybe necessary to delete!
		}
	}
}

func (dynamicSupervisor *DynamicSupervisor) pushRecreatedWorkingActorToArray(recreatedActor *workeractor.Actor) {
	//log.Println("recreate actor with identity:", recreatedActor.ActorProps.Identity)

	actorPoolName := utils.GetActorPollNameByActorIdentity(recreatedActor.ActorProps.Identity)

	*actorregistry.MyActorRegistry.FindActorByName(actorPoolName).(*[]actorabstraction.IActor) =
		append(*actorregistry.MyActorRegistry.FindActorByName(actorPoolName).(*[]actorabstraction.IActor),
			recreatedActor)
}

func (dynamicSupervisor *DynamicSupervisor) sendToRouter(myChan chan interface{}) {
	//changes here
	/*go func() {
		for msg := range myChan {
			actorregistry.MyActorRegistry.FindActorByName("routerActor").(*routeractor.RouterActor).SendMessage(msg)
		}
	}()*/
	for msg := range myChan {
		actorregistry.MyActorRegistry.FindActorByName("routerActor").(*routeractor.RouterActor).SendMessage(msg)
	}
}

func (dynamicSupervisor *DynamicSupervisor) collectDataOfActorsToBeDroppedToSingleChan(numberOfActors int, poolName string) chan interface{} {
	var actorsChanToDelete []chan interface{}
	for i := 0; i < numberOfActors; i++ {
		actorsChanToDelete = append(actorsChanToDelete,
			(*actorregistry.MyActorRegistry.FindActorByName(poolName).(*[]actorabstraction.IActor))[i].(*workeractor.Actor).ActorProps.ChanToReceiveData)
	}
	return utils.MergeWaitInterface(actorsChanToDelete)
}
