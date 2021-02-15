package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/autoscaleractor"
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/actor-model/dynamicsupervisor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
)

func RunApp(arr []string) error {
	dynamicSupervisor := dynamicsupervisor.NewDynamicSupervisor("dynamic_supervisor")
	//my connection actor
	connectionMaker := connectionactor.NewConnectionActor("connection", arr)

	actorPoll, err := dynamicSupervisor.CreateActorPoll(5)
	if err != nil {
		return err
	}

	//my router actor
	routerActor := routeractor.NewRouterActor("router", actorPoll)

	autoscalingActor := autoscaleractor.NewAutoscalingActor("autoscaling", dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)

	connectionMaker.SendDataToMultipleActorsOverChan(routerActor.ChanToRecvMsg, autoscalingActor.ChanToReceiveMessagesForCount)

	return nil
}
