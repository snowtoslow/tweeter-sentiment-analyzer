package aggregatoractor

import (
	"fmt"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/clientactor"
	"tweeter-sentiment-analyzer/actor-model/sinkactor"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
)

func NewAggregatorActor(actorName string) actorabstraction.IActor {
	chanForMessages := make(chan interface{}, constants.GlobalChanSize)

	aggregatorActor := &AggregatorActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanForMessages,
		},
		StorageToAggregateTweets: map[string]interface{}{},
	}

	(*actorregistry.MyActorRegistry)["aggregatorActor"] = aggregatorActor

	go aggregatorActor.ActorLoop()

	return aggregatorActor
}

func (aggregatorActor *AggregatorActor) ActorLoop() {
	defer close(aggregatorActor.ActorProps.ChanToReceiveData)
	for {
		select {
		case action := <-aggregatorActor.ActorProps.ChanToReceiveData:
			if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct || fmt.Sprintf("%T", action) == constants.RetweetedStatus {
				aggregatorActor.extractAndSendUserByInterfaceType(action)
				aggregatorActor.pushIntoTweetStorage(action)
			} else {
				aggregatorActor.addTweetFields(action)
			}
		}
	}
}

func (aggregatorActor *AggregatorActor) addTweetFields(action interface{}) {
	var myVal interface{}
	if fmt.Sprintf("%T", action) == constants.EngagementRatio {
		if val, ok := aggregatorActor.StorageToAggregateTweets[action.(*models.EngagementRation).UniqueId]; ok {
			aggregatorActor.addEngRation(val, action.(*models.EngagementRation).Ratio)
			myVal = val
		}
	} else if fmt.Sprintf("%T", action) == constants.SentimentAnalysis {
		if val, ok := aggregatorActor.StorageToAggregateTweets[action.(*models.SentimentAnalysis).UniqueId]; ok {
			aggregatorActor.addSentimentAnalysis(val, action.(*models.SentimentAnalysis).Score)
			myVal = val
		}
	}
	actorregistry.MyActorRegistry.FindActorByName("sinkActor").(*sinkactor.SinkActor).SendMessage(myVal)
	actorregistry.MyActorRegistry.FindActorByName("clientActor").(*clientactor.ClientActor).SendMessage(myVal)
	delete(aggregatorActor.StorageToAggregateTweets, aggregatorActor.getIdByInterfaceType(myVal))
}

func (aggregatorActor *AggregatorActor) extractAndSendUserByInterfaceType(value interface{}) {
	if fmt.Sprintf("%T", value) == constants.JsonNameOfStruct {
		actorregistry.MyActorRegistry.FindActorByName("sinkActor").(*sinkactor.SinkActor).SendMessage(value.(*models.MyJsonName).Message.Tweet.User)
		actorregistry.MyActorRegistry.FindActorByName("clientActor").(*clientactor.ClientActor).SendMessage(value.(*models.MyJsonName).Message.Tweet.User) //send user to client
	} else if fmt.Sprintf("%T", value) == constants.RetweetedStatus {
		actorregistry.MyActorRegistry.FindActorByName("sinkActor").(*sinkactor.SinkActor).SendMessage(value.(models.RetweetedStatus).User)
		actorregistry.MyActorRegistry.FindActorByName("clientActor").(*clientactor.ClientActor).SendMessage(value.(models.RetweetedStatus).User) //send user to client
	}
}

func (aggregatorActor *AggregatorActor) getIdByInterfaceType(value interface{}) (keyId string) {
	if fmt.Sprintf("%T", value) == constants.RetweetedStatusPointer {
		keyId = value.(*models.RetweetedStatus).UniqueId
	} else if fmt.Sprintf("%T", value) == constants.JsonNameOfStruct {
		keyId = value.(*models.MyJsonName).Message.UniqueId
	}
	return
}

func (aggregatorActor *AggregatorActor) addEngRation(val interface{}, engRatio float64) {
	if fmt.Sprintf("%T", val) == constants.RetweetedStatusPointer {
		val.(*models.RetweetedStatus).AggregationRation = engRatio
	} else if fmt.Sprintf("%T", val) == constants.JsonNameOfStruct {
		val.(*models.MyJsonName).Message.AggregationRation = engRatio
	}
}

func (aggregatorActor *AggregatorActor) addSentimentAnalysis(val interface{}, sentScore int8) {
	if fmt.Sprintf("%T", val) == constants.RetweetedStatusPointer {
		val.(*models.RetweetedStatus).SentimentSCore = sentScore
	} else if fmt.Sprintf("%T", val) == constants.JsonNameOfStruct {
		val.(*models.MyJsonName).Message.SentimentSCore = sentScore
	}
}

func (aggregatorActor *AggregatorActor) pushIntoTweetStorage(action interface{}) {
	if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct {
		aggregatorActor.StorageToAggregateTweets[action.(*models.MyJsonName).Message.UniqueId] = action.(*models.MyJsonName)
	} else if fmt.Sprintf("%T", action) == constants.RetweetedStatus {
		aggregatorActor.StorageToAggregateTweets[action.(models.RetweetedStatus).UniqueId] = action.(models.RetweetedStatus)
	}
}

func (aggregatorActor *AggregatorActor) SendMessage(msg interface{}) {
	aggregatorActor.ActorProps.ChanToReceiveData <- msg
}
