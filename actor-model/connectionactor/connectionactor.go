package connectionactor

import "tweeter-sentiment-analyzer/actor-model/actorabstraction"

type ConnectionActor struct {
	/*Identity       string
	ChanToSendData chan string*/
	ActorProps actorabstraction.AbstractActor
	Routes     []string
}
