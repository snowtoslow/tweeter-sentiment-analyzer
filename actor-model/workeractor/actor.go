package workeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type Actor struct {
	/*Identity          string
	ChanToReceiveData chan string*/
	ActorProps actorabstraction.AbstractActor
}
