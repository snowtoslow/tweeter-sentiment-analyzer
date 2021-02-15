package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/autoscaleractor"
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/actor-model/dynamicsupervisor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
)

func RunApp(arr []string) error {
	dynamicSupervisor := dynamicsupervisor.NewDynamicSupervisor("dynamic_supervisor")
	//my connection workeractor
	connectionMaker := connectionactor.NewConnectionActor("connection")

	actorPoll, err := dynamicSupervisor.CreateActorPoll(5)
	if err != nil {
		return err
	}

	//my router workeractor
	routerActor := routeractor.NewRouterActor("router", actorPoll)

	autoscalingActor := autoscaleractor.NewAutoscalingActor("autoscaling", dynamicSupervisor.ChanToReceiveNumberOfActorsToCreate)

	connectionMaker.SendDataToMultipleActorsOverChan(arr, routerActor.ChanToRecvMsg, autoscalingActor.ChanToReceiveMessagesForCount)

	//actorregistry.MyActorRegistry.TestFindActorByName("dynamicSupervisor")

	return nil
}
