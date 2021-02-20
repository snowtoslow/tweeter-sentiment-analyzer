package dynamicsupervisor

import "tweeter-sentiment-analyzer/actor-model/actorabstraction"

type DynamicSupervisor struct {
	ActorProps                          actorabstraction.AbstractActor
	ChanToReceiveNumberOfActorsToCreate chan int
}
