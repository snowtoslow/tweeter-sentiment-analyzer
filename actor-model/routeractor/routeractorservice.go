package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/constants"
)

func NewRouterActor(actorName string) actorabstraction.IActor {
	chanToRecvMsg := make(chan interface{}, constants.GlobalChanSize)

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

func (routerActor *RouterActor) SendMessage(data interface{}) {
	routerActor.ActorProps.ChanToReceiveData <- data
}

func (routerActor *RouterActor) ActorLoop() {
	defer close(routerActor.ActorProps.ChanToReceiveData)
	for {
		select {
		case output := <-routerActor.ActorProps.ChanToReceiveData:
			if routerActor.CurrentActorIndex >= len(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor)) {
				routerActor.CurrentActorIndex = 0
			}
			(*actorregistry.MyActorRegistry.FindActorByName("actorPool").(*[]actorabstraction.IActor))[routerActor.CurrentActorIndex].SendMessage(output) // change here from routerActor.Actors[routerActor.CurrentActorIndex].ChanToReceiveData <- output
			routerActor.CurrentActorIndex++
		}
	}
}
