package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type RouterActor struct {
	ActorProps        actorabstraction.AbstractActor
	CurrentActorIndex int
}
