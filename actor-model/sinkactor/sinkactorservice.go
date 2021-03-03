package sinkactor

import (
	"fmt"
	"log"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/constants"
)

func NewSinkActor(actorName string) actorabstraction.IActor {
	chanToRecv := make(chan interface{}, constants.GlobalChanSize)
	sinkActor := &SinkActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToRecv,
		},
	}

	go sinkActor.ActorLoop()

	(*actorregistry.MyActorRegistry)["sinkActor"] = sinkActor

	return sinkActor
}

func (sinkActor *SinkActor) ActorLoop() {
	defer close(sinkActor.ActorProps.ChanToReceiveData)
	for {
		action := <-sinkActor.ActorProps.ChanToReceiveData
		if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct {
			//log.Println("TWEET ID:",action.(*models.MyJsonName).Message.UniqueId)
		} else {
			log.Printf("STRUCT WITH TWEET ID: %T", action)
		}
		//log.Println("SINK ACTOR:",<-sinkActor.ActorProps.ChanToReceiveData)
	}
}

func (sinkActor *SinkActor) SendMessage(msg interface{}) {
	sinkActor.ActorProps.ChanToReceiveData <- msg
}
