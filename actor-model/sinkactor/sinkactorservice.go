package sinkactor

import (
	"fmt"
	"log"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
)

func NewSinkActor(actorName string) actorabstraction.IActor {
	chanToRecv := make(chan interface{}, constants.GlobalChanSize)
	sinkActor := &SinkActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToRecv,
		},
		TweetStorage: map[string]interface{}{},
	}

	go sinkActor.ActorLoop()

	(*actorregistry.MyActorRegistry)["sinkActor"] = sinkActor

	return sinkActor
}

func (sinkActor *SinkActor) ActorLoop() {
	defer close(sinkActor.ActorProps.ChanToReceiveData)
	for {
		select {
		case action := <-sinkActor.ActorProps.ChanToReceiveData:
			sinkActor.SinkBuffer = append(sinkActor.SinkBuffer, action)
			if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct || fmt.Sprintf("%T", action) == constants.RetweetedStatus {
				//log.Println("TWEET ID:",action.(*models.MyJsonName).Message.UniqueId)
				sinkActor.pushIntoTweetStorage(action)
			}

			if len(sinkActor.SinkBuffer) == 128 {
				log.Println("CLEAR BUFFER")
				//log.Println("SINK:",sinkActor.SinkBuffer)
				sinkActor.getSmthTest()
				sinkActor.SinkBuffer = sinkActor.SinkBuffer[:0]
			}
		}
	}
}

func (sinkActor *SinkActor) getSmthTest() {
	for _, v := range sinkActor.SinkBuffer {
		if fmt.Sprintf("%T", v) != constants.JsonNameOfStruct && fmt.Sprintf("%T", v) != constants.RetweetedStatus {
			switch v.(type) {
			case *models.EngagementRation:
				//log.Println("ENG RATIO:")
				if val, ok := sinkActor.TweetStorage[v.(*models.EngagementRation).UniqueId]; ok {
					log.Println("found in eng ration!", len(sinkActor.TweetStorage))
					sinkActor.addEngRation(val, v.(*models.EngagementRation).Ratio)
					delete(sinkActor.TweetStorage, v.(*models.EngagementRation).UniqueId)
					//log.Println("after delete:",len(sinkActor.TweetStorage))
				} else {
					log.Println("not found in eng ratio!")
				}
			case *models.SentimentAnalysis:
				//log.Println("SENT ANALYSIS:")
				if val, ok := sinkActor.TweetStorage[v.(*models.SentimentAnalysis).UniqueId]; ok {
					log.Println("found in sent anal!", len(sinkActor.TweetStorage))
					sinkActor.addSentimentAnalysis(val, v.(*models.SentimentAnalysis).Score)
					delete(sinkActor.TweetStorage, v.(*models.SentimentAnalysis).UniqueId)
					//log.Println("after delete:",len(sinkActor.TweetStorage))
				} else {
					log.Println("not found in eng analysis")
				}
			}

		}

	}
}

func (sinkActor *SinkActor) addEngRation(val interface{}, engRatio float64) {
	switch val.(type) {
	case *models.RetweetedStatus:
		val.(*models.RetweetedStatus).AggregationRation = engRatio
	case *models.MyJsonName:
		val.(*models.MyJsonName).Message.AggregationRation = engRatio
	}
}

func (sinkActor *SinkActor) addSentimentAnalysis(val interface{}, sentScore int8) {
	switch val.(type) {
	case *models.RetweetedStatus:
		val.(*models.RetweetedStatus).SentimentSCore = sentScore
	case *models.MyJsonName:
		val.(*models.MyJsonName).Message.SentimentSCore = sentScore
	}
}

func (sinkActor *SinkActor) pushIntoTweetStorage(action interface{}) {
	if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct {
		//log.Println("TWEET ID:",action.(*models.MyJsonName).Message.UniqueId)
		sinkActor.TweetStorage[action.(*models.MyJsonName).Message.UniqueId] = action.(*models.MyJsonName)
	} else if fmt.Sprintf("%T", action) == constants.RetweetedStatus {
		sinkActor.TweetStorage[action.(models.RetweetedStatus).UniqueId] = action.(models.RetweetedStatus)
	}
}

func (sinkActor *SinkActor) SendMessage(msg interface{}) {
	sinkActor.ActorProps.ChanToReceiveData <- msg
}
