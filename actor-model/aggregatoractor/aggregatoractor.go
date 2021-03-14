package aggregatoractor

import "tweeter-sentiment-analyzer/actor-model/actorabstraction"

type AggregatorActor struct {
	ActorProps               actorabstraction.AbstractActor
	StorageToAggregateTweets map[string]interface{}
}
