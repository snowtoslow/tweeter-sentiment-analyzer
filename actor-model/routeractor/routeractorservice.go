package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/workeractor"
	"tweeter-sentiment-analyzer/constants"
)

func NewRouterActor(actorName string) *RouterActor {
	chanToRecvMsg := make(chan string, constants.GlobalChanSize)

	routerActor := &RouterActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToRecvMsg,
		},
		CurrentActorIndex: 0,
	}

	go routerActor.ActorLoop()

	(*actorregistry.MyActorRegistry)["routerActor"] = routerActor

	return routerActor
}

func (routerActor *RouterActor) SendMessage(data string) {
	routerActor.ActorProps.ChanToReceiveData <- data
}

func (routerActor *RouterActor) ActorLoop() {
	defer close(routerActor.ActorProps.ChanToReceiveData)
	actors := actorregistry.MyActorRegistry.TestFindActorByName("actorPool").([]workeractor.Actor)
	for {
		select {
		case output := <-routerActor.ActorProps.ChanToReceiveData:
			if routerActor.CurrentActorIndex >= len(actors) {
				routerActor.CurrentActorIndex = 0
			}
			(actors)[routerActor.CurrentActorIndex].SendMessage(output) // change here from routerActor.Actors[routerActor.CurrentActorIndex].ChanToReceiveData <- output
			routerActor.CurrentActorIndex++
		}
	}
}
