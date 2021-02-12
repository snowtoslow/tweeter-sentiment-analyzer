package autoscaleractor

import (
	"log"
	"regexp"
	"time"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
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
	prevCounter := 0
	for {
		select {
		case msg := <-autoscalingActor.ChanToReceiveMessagesForCount:
			action := regexp.MustCompile(constants.JsonRegex).FindString(msg)
			if len(action) != 0 {
				counter++
			}
		case <-ticker.C:
			//log.Println("COUNTER:", counter)
			//log.Println("PREV COUNTER:",prevCounter)
			prevCounter = counter
			movingAverage := utils.MovingExpAvg(float64(counter), float64(prevCounter), 1, 2)
			log.Println("Actor number:", int(movingAverage/20))
			//autoscalingActor.ChanToSendCounterResult<-int(movingAverage/20)
			counter = 0
		}
	}
}
