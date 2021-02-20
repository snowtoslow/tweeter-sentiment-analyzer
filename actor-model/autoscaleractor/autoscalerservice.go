package autoscaleractor

import (
	"regexp"
	"time"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/dynamicsupervisor"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewAutoscalingActor(actorName string) actorabstraction.IActor {
	chanForMessages := make(chan interface{}, constants.GlobalChanSize)

	autoscalingActor := &AutoscalingActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanForMessages,
		},
	}

	(*actorregistry.MyActorRegistry)["autoscalingActor"] = autoscalingActor

	go autoscalingActor.ActorLoop()

	return autoscalingActor
}

func (autoscalingActor *AutoscalingActor) ActorLoop() {
	defer close(autoscalingActor.ActorProps.ChanToReceiveData)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	counter := 0
	prevCounter := 0
	prevMovingAvg := 0
	movingAvg := 0
	for {
		select {
		case msg := <-autoscalingActor.ActorProps.ChanToReceiveData:
			action := regexp.MustCompile(constants.JsonRegex).FindString(msg.(string))
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
			//log.Println("HERE:",actorregistry.MyActorRegistry.FindActorByName("dynamic_supervisor"))
			autoscalingActor.sendMessageToSupervisor(movingAvg - prevMovingAvg)
			//log.Println("prev mov avg:",prevMovingAvg)
			//log.Println("moving avg:",movingAvg)
		}
	}
}

func (autoscalingActor *AutoscalingActor) SendMessage(msg interface{}) {
	autoscalingActor.ActorProps.ChanToReceiveData <- msg
}

func (autoscalingActor *AutoscalingActor) sendMessageToSupervisor(msg interface{}) {
	actorregistry.MyActorRegistry.FindActorByName("dynamicSupervisor").(*dynamicsupervisor.DynamicSupervisor).SendMessage(msg)
}
