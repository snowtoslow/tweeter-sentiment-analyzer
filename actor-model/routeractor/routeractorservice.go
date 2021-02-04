package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actor"
)

func NewRouterActor(actorName string, actorAmount int) (*RouterActor, error) {
	chanToRecvMsg := make(chan string, 10)
	actorPool, err := actor.CreateActorPoll(actorAmount) // actor pool created here!
	if err != nil {
		return nil, err
	}
	routerActor := &RouterActor{
		Identity:          actorName + "_actor",
		ChanToRecvMsg:     chanToRecvMsg,
		CurrentActorIndex: 0,
		Actors:            actorPool,
	}

	go routerActor.actorLoop() //actor loop for balancing;

	return routerActor, nil
}

func (routerActor *RouterActor) SendMessage(data string) {
	routerActor.ChanToRecvMsg <- data
}

func (routerActor *RouterActor) actorLoop() {
	defer close(routerActor.ChanToRecvMsg)
	for {
		select {
		case output := <-routerActor.ChanToRecvMsg:
			if routerActor.CurrentActorIndex >= len(routerActor.Actors) {
				routerActor.CurrentActorIndex = 0
			}
			routerActor.Actors[routerActor.CurrentActorIndex].ChanToReceiveData <- output
			routerActor.CurrentActorIndex++
		}
	}
}
