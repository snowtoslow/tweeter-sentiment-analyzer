package routeractor

import (
	"log"
	"tweeter-sentiment-analyzer/actor-model/actor"
)

func NewRouterActor(actorName string) *RouterActor {
	chanToRecvMsg := make(chan string, 10)
	routerActor := &RouterActor{
		Identity:      actorName + "_actor",
		ChanToRecvMsg: chanToRecvMsg,
	}

	return routerActor
}

func (routerActor *RouterActor) SendMessage(data string, randomActor *actor.Actor) {
	routerActor.ChanToRecvMsg <- data
	action := <-routerActor.ChanToRecvMsg
	randomActor.ChanToReceiveData <- action
	log.Printf("DATA WAS SENT FROM MAIN ROUTER ACTOR:%s to WORKER_ACTOR: %s", routerActor.Identity, randomActor.Identity)
}
