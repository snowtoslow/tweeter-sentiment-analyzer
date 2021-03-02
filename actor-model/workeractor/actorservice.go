package workeractor

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	message_types "tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
	"tweeter-sentiment-analyzer/utils"
)

func NewActor(actorName string, dynamicSup IDynamicSupervisor) actorabstraction.IActor {
	chanToRecv := make(chan interface{}, constants.GlobalChanSize)
	actor := &Actor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName,
			ChanToReceiveData: chanToRecv,
		},
		DynamicSupervisorAvoidance: dynamicSup,
	}

	go actor.ActorLoop()

	return actor
}

func (actor *Actor) ActorLoop() {
	defer close(actor.ActorProps.ChanToReceiveData)
	for {
		action := actor.processReceivedMessage(<-actor.ActorProps.ChanToReceiveData)
		if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct {
			actor.delegateWork(action)
		} else if fmt.Sprintf("%T", action) == constants.PanicMessageType {
			log.Println("ERROR:")
			errMsg := message_types.ErrorToSupervisor{
				ActorIdentity: actor.ActorProps.Identity,
				Message:       message_types.PanicMessage("error occurred in worker actor with identity " + actor.ActorProps.Identity),
			}
			actor.SendMessageToSupervisor(errMsg)
		}
	}
}

func (actor *Actor) delegateWork(action interface{}) {
	if strings.Contains(actor.ActorProps.Identity, constants.SentimentActorPool) {
		//log.Printf("TEXT:%s\nRESULT:%v\n", action.(*models.MyJsonName).Message.Tweet.Text, utils.AnalyzeSentiments(action.(*models.MyJsonName).Message.Tweet.Text))
		//log.Println("SENTIMENTS ANALySYS")
	} else if strings.Contains(actor.ActorProps.Identity, constants.AggregationActorPool) {
		log.Println("ENGAGEMENT RATIO:", utils.EngagementRatio(action.(*models.MyJsonName).Message.Tweet.RetweetedStatus,
			action.(*models.MyJsonName).Message.Tweet.User.FavouritesCount,
			action.(*models.MyJsonName).Message.Tweet.User.FollowersCount))
		//log.Println("AGGREGATION ACTOR")
	}
}

func (actor *Actor) SendMessage(data interface{}) {
	actor.ActorProps.ChanToReceiveData <- data
}

func (actor *Actor) SendMessageToSupervisor(msg interface{}) {
	actor.DynamicSupervisorAvoidance.SendErrMessage(msg)
}

func (actor *Actor) processReceivedMessage(dataToRecv interface{}) (resultDataStructure interface{}) {
	regexData := regexp.MustCompile(constants.JsonRegex)
	if receivedString := regexData.FindString(dataToRecv.(string)); len(receivedString) != 0 {
		// be carefully with error from json CreateMessageType
		resultDataStructure = utils.CreateMessageType(receivedString)
	}
	return
}
