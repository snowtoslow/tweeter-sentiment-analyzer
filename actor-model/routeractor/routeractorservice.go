package routeractor

import (
	"bufio"
	"net/http"
	"tweeter-sentiment-analyzer/actor-model/actor"
	"tweeter-sentiment-analyzer/constants"
)

func NewRouterActor(actorName string, actorAmount int) (*RouterActor, error) {
	chanToRecvMsg := make(chan string, constants.GlobalChanSize)
	actorPool, err := actor.CreateActorPoll(actorAmount) // actor pool created here!
	if err != nil {
		return nil, err
	}
	routerActor := &RouterActor{
		Identity:          actorName + constants.ActorName,
		ChanToRecvMsg:     chanToRecvMsg,
		CurrentActorIndex: 0,
		Actors:            actorPool,
	}

	go routerActor.actorLoop() //actor loop for balancing;

	return routerActor, nil
}

func (routerActor *RouterActor) MakeRequest(url string, ch chan string) {
	res, err := http.Get(url)
	if err != nil {
		close(ch)
		return
	}
	defer res.Body.Close()
	defer close(ch)
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		routerActor.sendMessage(scanner.Text())
	}
}

func (routerActor *RouterActor) sendMessage(data string) {
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
			routerActor.Actors[routerActor.CurrentActorIndex].SendMessage(output) // change here from routerActor.Actors[routerActor.CurrentActorIndex].ChanToReceiveData <- output
			routerActor.CurrentActorIndex++
		}
	}
}
