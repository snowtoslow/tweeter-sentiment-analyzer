package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/actor"
	"tweeter-sentiment-analyzer/actor-model/autoscaleractor"
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
)

func RunApp(arr []string) error {
	//my connection actor
	connectionMaker := connectionactor.NewConnectionActor("connection", arr)

	actorPoll, err := actor.CreateActorPoll(5)
	if err != nil {
		return err
	}

	//my router actor
	routerActor := routeractor.NewRouterActor("router", actorPoll)

	autoscalingActor := autoscaleractor.NewAutoscalingActor("autoscaller")

	connectionMaker.SendDataToMultipleActorsOverChan(routerActor.ChanToRecvMsg, autoscalingActor.ChanToReceiveMessagesForCount)

	return nil
}
