package supervisor

import (
	"fmt"
	"log"
	"tweeter-sentiment-analyzer/actor-model/actor"
	"tweeter-sentiment-analyzer/actor-model/messagetypes"
	"tweeter-sentiment-analyzer/constants"
)

func NewSupervisor(actors []*actor.Actor) *Supervisor {
	ch := make(chan messagetypes.ErrorForSupervisor, constants.GlobalChanSize)

	supervisor := &Supervisor{
		ChanToReceiveErrorMessage: ch,
		Actors:                    actors,
	}

	go supervisor.actorLoop()

	return supervisor
}

func (supervisor *Supervisor) SendMessage(errMsg messagetypes.ErrorForSupervisor) {
	supervisor.ChanToReceiveErrorMessage <- errMsg
}

func (supervisor *Supervisor) actorLoop() {
	log.Println("actor loop!")
}

func (supervisor *Supervisor) deleteActorByActorName() (actors []*actor.Actor, err error) {
	crashedActorIdentity := <-supervisor.ChanToReceiveErrorMessage
	for k, v := range supervisor.Actors {
		if v.Identity == crashedActorIdentity.FailedActorIdentity {
			return append(supervisor.Actors[:k], supervisor.Actors[k+1:]...), nil
		}
	}
	return nil, fmt.Errorf("actor with this identity: %s does not exist", crashedActorIdentity.FailedActorIdentity)
}
