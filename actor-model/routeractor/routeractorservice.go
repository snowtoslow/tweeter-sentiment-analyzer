package routeractor

import (
	worker_actor "tweeter-sentiment-analyzer/actor-model/workeractor"
	"tweeter-sentiment-analyzer/constants"
)

func NewRouterActor(actorName string, actorPoll *[]worker_actor.Actor) *RouterActor {
	chanToRecvMsg := make(chan string, constants.GlobalChanSize)

	routerActor := &RouterActor{
		Identity:          actorName + constants.ActorName,
		ChanToRecvMsg:     chanToRecvMsg,
		CurrentActorIndex: 0,
		Actors:            actorPoll,
	}

	go routerActor.actorLoop() //workeractor loop for balancing;

	return routerActor
}

func (routerActor *RouterActor) SendMessage(data string) {
	routerActor.ChanToRecvMsg <- data
}

func (routerActor *RouterActor) actorLoop() {
	defer close(routerActor.ChanToRecvMsg)
	for {
		select {
		case output := <-routerActor.ChanToRecvMsg:
			if routerActor.CurrentActorIndex >= len(*routerActor.Actors) {
				routerActor.CurrentActorIndex = 0
			}
			(*routerActor.Actors)[routerActor.CurrentActorIndex].SendMessage(output) // change here from routerActor.Actors[routerActor.CurrentActorIndex].ChanToReceiveData <- output
			routerActor.CurrentActorIndex++
		}
	}
}
