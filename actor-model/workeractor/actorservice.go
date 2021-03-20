package workeractor

import (
	"fmt"
	"regexp"
	"strings"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/aggregatoractor"
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
			if &action.(*models.MyJsonName).Message.Tweet.RetweetedStatus != nil {
				actor.extractSubTweetsAndAnalyze(action)
			}
			generatedId := utils.GenerateUuidgen()
			action.(*models.MyJsonName).Message.UniqueId = generatedId
			//add unique id to user from extracted tweet:
			action.(*models.MyJsonName).Message.Tweet.User.UniqueId = generatedId
			actorregistry.MyActorRegistry.FindActorByName("aggregatorActor").(*aggregatoractor.AggregatorActor).SendMessage(action)
			actor.delegateWork(action.(*models.MyJsonName).Message.Tweet.Text,
				action.(*models.MyJsonName).Message.Tweet.RetweetedStatus,
				action.(*models.MyJsonName).Message.Tweet.User.FavouritesCount,
				action.(*models.MyJsonName).Message.Tweet.User.FollowersCount, generatedId)
		} else if fmt.Sprintf("%T", action) == constants.PanicMessageType {
			//log.Println("ERROR:")
			errMsg := message_types.ErrorToSupervisor{
				ActorIdentity: actor.ActorProps.Identity,
				Message:       message_types.PanicMessage("error occurred in worker actor with identity " + actor.ActorProps.Identity),
			}
			actor.SendMessageToSupervisor(errMsg)
		}
	}
}

func (actor *Actor) extractSubTweetsAndAnalyze(mainTweet interface{}) {
	generatedId := utils.GenerateUuidgen()
	mainTweet.(*models.MyJsonName).Message.Tweet.RetweetedStatus.UniqueId = generatedId
	//add unique id to user from extracted tweet:
	mainTweet.(*models.MyJsonName).Message.Tweet.User.UniqueId = generatedId
	actorregistry.MyActorRegistry.FindActorByName("aggregatorActor").(*aggregatoractor.AggregatorActor).SendMessage(mainTweet.(*models.MyJsonName).Message.Tweet.RetweetedStatus)
	actor.delegateWork(mainTweet.(*models.MyJsonName).Message.Tweet.RetweetedStatus.Text,
		mainTweet.(*models.MyJsonName).Message.Tweet.RetweetedStatus,
		mainTweet.(*models.MyJsonName).Message.Tweet.RetweetedStatus.FavoriteCount,
		mainTweet.(*models.MyJsonName).Message.Tweet.RetweetedStatus.User.FollowersCount,
		generatedId)
}

func (actor *Actor) delegateWork(textForSentimentAnalysis string, retweetedStatus models.RetweetedStatus, favCount int64, followersCount int64, generatedId string) {
	if strings.Contains(actor.ActorProps.Identity, constants.SentimentActorPool) {
		actorregistry.MyActorRegistry.FindActorByName("aggregatorActor").(*aggregatoractor.AggregatorActor).SendMessage(&models.SentimentAnalysis{
			Score:    utils.AnalyzeSentiments(textForSentimentAnalysis),
			UniqueId: generatedId,
		})
	} else if strings.Contains(actor.ActorProps.Identity, constants.AggregationActorPool) {
		actorregistry.MyActorRegistry.FindActorByName("aggregatorActor").(*aggregatoractor.AggregatorActor).SendMessage(&models.EngagementRation{
			Ratio:    utils.EngagementRatio(retweetedStatus, favCount, followersCount),
			UniqueId: generatedId,
		})
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
