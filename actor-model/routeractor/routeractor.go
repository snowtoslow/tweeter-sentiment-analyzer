package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/routerstrategy"
)

type RouterActor struct {
	ActorProps      actorabstraction.AbstractActor
	RoutingStrategy *routerstrategy.RoundRobin
}
