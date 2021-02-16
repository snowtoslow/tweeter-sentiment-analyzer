package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/workeractor"
	"tweeter-sentiment-analyzer/constants"
)

func NewRouterActor(actorName string) *RouterActor {
	chanToRecvMsg := make(chan string, constants.GlobalChanSize)

	routerActor := &RouterActor{
		Identity:          actorName + constants.ActorName,
		ChanToRecvMsg:     chanToRecvMsg,
		CurrentActorIndex: 0,
	}

	(*actorregistry.MyActorRegistry)["routerActor"] = routerActor

	go routerActor.ActorLoop() //workeractor loop for balancing;

	return routerActor
}

func (routerActor *RouterActor) SendMessage(data string) {
	routerActor.ChanToRecvMsg <- data
}

func (routerActor *RouterActor) ActorLoop() {
	defer close(routerActor.ChanToRecvMsg)
	actors := actorregistry.MyActorRegistry.TestFindActorByName("actorPool").([]workeractor.Actor)
	for {
		select {
		case output := <-routerActor.ChanToRecvMsg:
			if routerActor.CurrentActorIndex >= len(actors) {
				routerActor.CurrentActorIndex = 0
			}
			(actors)[routerActor.CurrentActorIndex].SendMessage(output) // change here from routerActor.Actors[routerActor.CurrentActorIndex].ChanToReceiveData <- output
			routerActor.CurrentActorIndex++
		}
	}
}
