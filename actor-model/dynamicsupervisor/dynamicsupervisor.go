package dynamicsupervisor

import "tweeter-sentiment-analyzer/actor-model/actorabstraction"

type DynamicSupervisor struct {
	/*Identity                            string
	ChanToReceiveErrors                 chan string*/
	ActorProps                          actorabstraction.AbstractActor
	ChanToReceiveNumberOfActorsToCreate chan int
}
