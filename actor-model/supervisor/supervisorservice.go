package supervisor

import (
	"tweeter-sentiment-analyzer/actor-model/messagetypes"
	"tweeter-sentiment-analyzer/constants"
)

func NewSupervisor(actorName string) *Supervisor {
	ch := make(chan *messagetypes.ErrorForSupervisor, constants.GlobalChanSize)

	supervisor := &Supervisor{
		Identity:                  actorName + constants.ActorName,
		ChanToReceiveErrorMessage: ch,
	}

	go supervisor.actorLoop()

	return supervisor
}

func (supervisor *Supervisor) SendMessage(errMsg *messagetypes.ErrorForSupervisor) {
	supervisor.ChanToReceiveErrorMessage <- errMsg
}

func (supervisor *Supervisor) actorLoop() {
	defer close(supervisor.ChanToReceiveErrorMessage)
	for {
		action := <-supervisor.ChanToReceiveErrorMessage
		action.PanicWithRecoveryFunction()
	}
}
