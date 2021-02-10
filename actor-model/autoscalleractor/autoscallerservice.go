package autoscalleractor

import (
	"log"
	"time"
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
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	counter := 0
	for {
		select {
		case <-autoscaller.ChanToReceiveMessagesForCount:
			//log.Println(<-autoscaller.ChanToReceiveMessagesForCount)
		case <-ticker.C:
			counter += len(<-autoscaller.ChanToReceiveMessagesForCount)
			log.Println("COUNTER:", counter)
		}
	}

}
