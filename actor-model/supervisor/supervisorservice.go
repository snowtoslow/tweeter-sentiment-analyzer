package supervisor

import "tweeter-sentiment-analyzer/actor-model/actor"

func (supervisor *Supervisor) HandleError() {

}

func (supervisor *Supervisor) CreateActor() *actor.Actor {
	return &actor.Actor{}
}
