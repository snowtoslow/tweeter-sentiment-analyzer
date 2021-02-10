package autoscalleractor

import (
	"log"
	"tweeter-sentiment-analyzer/constants"
)

func NewAutoscallerActor(actorName string) *AutoscallerActor {
	chanForMessages := make(chan string, constants.GlobalChanSize)

	autoscallerActor := &AutoscallerActor{
		Identity:                      actorName + constants.ActorName,
		ChanToReceiveMessagesForCount: chanForMessages,
	}

	go autoscallerActor.actorLoop()

	return autoscallerActor
}

func (autoscaller *AutoscallerActor) sendMessage(data string) {
	autoscaller.ChanToReceiveMessagesForCount <- data
}

func (autoscaller *AutoscallerActor) actorLoop() {
	defer close(autoscaller.ChanToReceiveMessagesForCount)
	for {
		log.Println("AUTOSCALLER TEST:", <-autoscaller.ChanToReceiveMessagesForCount)
	}
}

/*func (autoscaller *AutoscallerActor) processReceivedMsg(dataToRecv string){
	regexData := regexp.MustCompile(constants.JsonRegex)
	receivedString := regexData.FindString(dataToRecv)

}*/
