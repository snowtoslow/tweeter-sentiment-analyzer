package dynamicsupervisor

import "tweeter-sentiment-analyzer/constants"

func NewDynamicSupervisor(actorName string) *DynamicSupervisor {
	chanToReceiveAmountOfActorsToCreate := make(chan int, 10)

	dynamicSupervisor := &DynamicSupervisor{
		Identity:                            "dynamic_supervisor" + constants.ActorName,
		ChanToReceiveNumberOfActorsToCreate: chanToReceiveAmountOfActorsToCreate,
	}

	return dynamicSupervisor
}

func (dynamicSupervisor *DynamicSupervisor) sendMessage(numberOfActorsToReceive int) {
	dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate <- numberOfActorsToReceive
}

func (dynamicSupervisor *DynamicSupervisor) actorLoop() {
	defer close(dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)
}

func (dynamicSupervisor *DynamicSupervisor) addActors() {

}

func (dynamicSupervisor *DynamicSupervisor) deleteActors() {

}

func (dynamicSupervisor *DynamicSupervisor) createMainActorPool() {}
