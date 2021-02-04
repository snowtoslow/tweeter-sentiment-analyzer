package routeractor

import (
	"encoding/json"
	"log"
	"regexp"
	"tweeter-sentiment-analyzer/actor-model/actor"
	msgType "tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
)

func NewRouterActor(actorName string) *RouterActor {
	chanToRecvMsg := make(chan string, 10)
	routerActor := &RouterActor{
		Identity:      actorName + "_actor",
		ChanToRecvMsg: chanToRecvMsg,
	}

	return routerActor
}

func (routerActor *RouterActor) SendProcessedMessage(data string, randomActor *actor.Actor) {
	/*routerActor.ChanToRecvMsg <- data
	action := <-routerActor.ChanToRecvMsg*/
	randomActor.ChanToReceiveData <- routerActor.getAndProcessMsg(data)
	log.Printf("DATA WAS SENT FROM MAIN ROUTER ACTOR:%s to WORKER_ACTOR: %s", routerActor.Identity, randomActor.Identity)
}

func (routerActor *RouterActor) getAndProcessMsg(data string) interface{} {
	regexData := regexp.MustCompile("\\{.*\\:\\{.*\\:.*\\}\\}|\\{(.*?)\\}") // already tested
	routerActor.ChanToRecvMsg <- data
	receivedString := regexData.FindString(<-routerActor.ChanToRecvMsg)
	var tweetMsg *models.MyJsonName
	if receivedString == constants.PanicMessage {
		return msgType.PanicMessage(receivedString)
	} else {
		err := json.Unmarshal([]byte(receivedString), &tweetMsg)
		if err != nil {
			return err
		}
		return tweetMsg
	}
	//return receivedString
	//return regexJson.FindString(regexData.FindString(<-routerActor.ChanToRecvMsg))
}
