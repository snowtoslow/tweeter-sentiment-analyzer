package autoscaleractor

import (
	"regexp"
	"time"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewAutoscalingActor(actorName string, ch chan int) *AutoscalingActor {
	chanForMessages := make(chan string, constants.GlobalChanSize)

	autoscalingActor := &AutoscalingActor{
		Identity:                      actorName + constants.ActorName,
		ChanToReceiveMessagesForCount: chanForMessages,
	}

	go autoscalingActor.actorLoop(ch)

	return autoscalingActor
}

func (autoscalingActor *AutoscalingActor) actorLoop(ch chan int) {
	defer close(autoscalingActor.ChanToReceiveMessagesForCount)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	counter := 0
	prevCounter := 0
	prevMovingAvg := 0
	movingAvg := 0
	for {
		select {
		case msg := <-autoscalingActor.ChanToReceiveMessagesForCount:
			action := regexp.MustCompile(constants.JsonRegex).FindString(msg)
			if len(action) != 0 {
				counter++
			}
			prevMovingAvg = movingAvg
		case <-ticker.C:
			//log.Println("COUNTER:", counter)
			//log.Println("PREV COUNTER:",prevCounter)
			movingAvg = int(utils.MovingExpAvg(float64(counter), float64(prevCounter), 1, 2)) / 15
			prevCounter = counter
			counter = 0
			autoscalingActor.sendMsgToSupervisor(movingAvg, prevMovingAvg, ch)
			//log.Println("prev mov avg:",prevMovingAvg)
			//log.Println("moving avg:",movingAvg)
		}
	}
}

func (autoscalingActor *AutoscalingActor) sendMsgToSupervisor(counter int, prevCounter int, supervisorChan chan int) {
	for counter := range autoscalingActor.sendCountedMessageToTmpChanTest(counter, prevCounter) {
		supervisorChan <- counter
	}
}

func (autoscalingActor *AutoscalingActor) sendCountedMessageToTmpChanTest(movAvg int, prevMovAvg int) chan int {
	autoscalingActor.ChanToSendCounterResult = make(chan int, constants.GlobalChanSize)
	go func() {
		autoscalingActor.ChanToSendCounterResult <- movAvg - prevMovAvg
		close(autoscalingActor.ChanToSendCounterResult)
	}()
	return autoscalingActor.ChanToSendCounterResult
}
