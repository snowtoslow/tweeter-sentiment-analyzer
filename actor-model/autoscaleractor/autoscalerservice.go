package autoscaleractor

import (
	"log"
	"time"
	"tweeter-sentiment-analyzer/constants"
)

func NewAutoscalingActor(actorName string) *AutoscalingActor {
	chanForMessages := make(chan string, constants.GlobalChanSize)

	autoscalingActor := &AutoscalingActor{
		Identity:                      actorName + constants.ActorName,
		ChanToReceiveMessagesForCount: chanForMessages,
	}

	go autoscalingActor.actorLoop()

	return autoscalingActor
}

func (autoscalingActor *AutoscalingActor) sendMessage(data string) {
	autoscalingActor.ChanToReceiveMessagesForCount <- data
}

func (autoscalingActor *AutoscalingActor) actorLoop() {
	defer close(autoscalingActor.ChanToReceiveMessagesForCount)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	counter := 0
	for {
		select {
		case msg := <-autoscalingActor.ChanToReceiveMessagesForCount:
			if len(msg) != 0 {
				counter++
			}
		case <-ticker.C:
			log.Println("Length:", counter)
			counter = 0
		}
	}
}
