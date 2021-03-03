package workeractor

import (
	"fmt"
	"regexp"
	"strings"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	message_types "tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/actor-model/sinkactor"
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
			generatedId := utils.GenerateUuidgen()
			action.(*models.MyJsonName).Message.UniqueId = generatedId
			actorregistry.MyActorRegistry.FindActorByName("sinkActor").(*sinkactor.SinkActor).SendMessage(action)
			actor.delegateWork(action, generatedId)
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

func (actor *Actor) delegateWork(action interface{}, generatedId string) {
	if strings.Contains(actor.ActorProps.Identity, constants.SentimentActorPool) {
		//log.Printf("TEXT:%s\nRESULT:%v\n", action.(*models.MyJsonName).Message.Tweet.Text, utils.AnalyzeSentiments(action.(*models.MyJsonName).Message.Tweet.Text))
		//log.Println("SENTIMENTS ANALySYS")
		actorregistry.MyActorRegistry.FindActorByName("sinkActor").(*sinkactor.SinkActor).SendMessage(struct {
			SentimentValue int8
			GeneratedId    string
		}{
			SentimentValue: utils.AnalyzeSentiments(action.(*models.MyJsonName).Message.Tweet.Text),
			GeneratedId:    generatedId,
		})
	} else if strings.Contains(actor.ActorProps.Identity, constants.AggregationActorPool) {
		/*log.Println("ENGAGEMENT RATIO:", utils.EngagementRatio(action.(*models.MyJsonName).Message.Tweet.RetweetedStatus,
		action.(*models.MyJsonName).Message.Tweet.User.FavouritesCount,
		action.(*models.MyJsonName).Message.Tweet.User.FollowersCount))*/
		//log.Println("AGGREGATION ACTOR")
		actorregistry.MyActorRegistry.FindActorByName("sinkActor").(*sinkactor.SinkActor).SendMessage(struct {
			EngagementRatio float64
			GeneratedId     string
		}{
			EngagementRatio: utils.EngagementRatio(action.(*models.MyJsonName).Message.Tweet.RetweetedStatus,
				action.(*models.MyJsonName).Message.Tweet.User.FavouritesCount,
				action.(*models.MyJsonName).Message.Tweet.User.FollowersCount),
			GeneratedId: generatedId,
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
