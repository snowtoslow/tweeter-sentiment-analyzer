package routeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

/*type RouterActor struct {
	Identity          string
	ChanToRecvMsg     chan string
	CurrentActorIndex int
}*/

type RouterActor struct {
	ActorProps        actorabstraction.AbstractActor
	CurrentActorIndex int
}
