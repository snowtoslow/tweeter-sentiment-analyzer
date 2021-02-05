package supervisor

import (
	"tweeter-sentiment-analyzer/actor-model/actor"
	message_types "tweeter-sentiment-analyzer/actor-model/messagetypes"
)

type Supervisor struct {
	Actors                    []*actor.Actor
	ChanToReceiveErrorMessage chan message_types.ErrorForSupervisor
}
