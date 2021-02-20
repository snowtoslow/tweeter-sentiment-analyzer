package workeractor

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type IDynamicSupervisor interface {
	SendErrMessage(msg interface{})
}

type Actor struct {
	DynamicSupervisorAvoidance IDynamicSupervisor
	ActorProps                 actorabstraction.AbstractActor
}
