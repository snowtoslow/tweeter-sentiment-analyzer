package connectionactor

import "tweeter-sentiment-analyzer/actor-model/actorabstraction"

type ConnectionActor struct {
	ActorProps actorabstraction.AbstractActor
	Routes     []string
}
