package autoscaleractor

import (
	"regexp"
	"time"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/dynamicsupervisor"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewAutoscalingActor(actorName string) *AutoscalingActor {
	chanForMessages := make(chan string, constants.GlobalChanSize)

	autoscalingActor := &AutoscalingActor{
		Identity:                      actorName + constants.ActorName,
		ChanToReceiveMessagesForCount: chanForMessages,
	}

	(*actorregistry.MyActorRegistry)["autoscalingActor"] = autoscalingActor

	go autoscalingActor.ActorLoop()

	return autoscalingActor
}

func (autoscalingActor *AutoscalingActor) ActorLoop() {
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
			//log.Println("HERE:",actorregistry.MyActorRegistry.TestFindActorByName("dynamic_supervisor"))
			autoscalingActor.SendMessage(movingAvg - prevMovingAvg)
			//log.Println("prev mov avg:",prevMovingAvg)
			//log.Println("moving avg:",movingAvg)
		}
	}
}

func (autoscalingActor *AutoscalingActor) SendMessage(msg int) {
	actorregistry.MyActorRegistry.TestFindActorByName("dynamicSupervisor").(*dynamicsupervisor.DynamicSupervisor).ChanToReceiveNumberOfActorsToCreate <- msg
}
