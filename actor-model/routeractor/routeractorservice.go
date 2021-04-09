package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/constants"
	"tweeter-sentiment-analyzer/actor-model/routerstrategy"
)

func NewRouterActor(actorName string) actorabstraction.IActor {
	chanToRecvMsg := make(chan interface{}, constants.GlobalChanSize)

	routerActor := &RouterActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToRecvMsg,
		},
		RoutingStrategy: routerstrategy.NewRoundRobinStrategy(),
	}

	go routerActor.ActorLoop()

	(*actorregistry.MyActorRegistry)["routerActor"] = routerActor

	return routerActor
}

func (routerActor *RouterActor) SendMessage(data interface{}) {
	routerActor.ActorProps.ChanToReceiveData <- data
}

func (routerActor *RouterActor) ReceiveMessageFromSupervisor(data interface{}) {
	//log.Println("RECEIVE DATA FROM SUPERVISOR:")
	routerActor.ActorProps.ChanToReceiveData <- data
}

func (routerActor *RouterActor) ActorLoop() {
	balancers := routerActor.RoutingStrategy.MultipleBalancerEntity(*actorregistry.MyActorRegistry.FindActorByName(constants.SentimentActorPool).(*[]actorabstraction.IActor), *actorregistry.MyActorRegistry.FindActorByName(constants.AggregationActorPool).(*[]actorabstraction.IActor))
	defer close(routerActor.ActorProps.ChanToReceiveData)
	for {
		select {
		case out := <-routerActor.ActorProps.ChanToReceiveData:
			for _, v := range balancers {
				v.Balancer(out)
			}
		}
	}
}
