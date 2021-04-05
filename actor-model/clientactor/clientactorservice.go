package clientactor

import (
	"log"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/constants"
)

func NewClientActor(actorName string) actorabstraction.IActor {
	chanForMessages := make(chan interface{}, constants.GlobalChanSize)

	clientActor := &ClientActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanForMessages,
		},
	}

	(*actorregistry.MyActorRegistry)["clientActor"] = clientActor

	go clientActor.ActorLoop()

	return clientActor
}

func (clientActor *ClientActor) ActorLoop() {
	defer close(clientActor.ActorProps.ChanToReceiveData)
	for {
		log.Printf("CLIENT ACTOR TYPE: %T", <-clientActor.ActorProps.ChanToReceiveData)
	}
}

func (clientActor *ClientActor) SendMessage(msg interface{}) {
	clientActor.ActorProps.ChanToReceiveData <- msg
}
