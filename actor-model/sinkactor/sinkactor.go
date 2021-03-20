package sinkactor

import "tweeter-sentiment-analyzer/actor-model/actorabstraction"

type SinkActor struct {
	ActorProps actorabstraction.AbstractActor
	sinkBuffer map[string][]interface{}
}
