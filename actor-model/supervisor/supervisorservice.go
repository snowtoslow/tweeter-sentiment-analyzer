package supervisor

import "tweeter-sentiment-analyzer/actor-model/workeractor"

func (supervisor *Supervisor) HandleError() {

}

func (supervisor *Supervisor) CreateActor() *workeractor.Actor {
	return &workeractor.Actor{}
}
