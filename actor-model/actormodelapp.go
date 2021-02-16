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
	connectionMaker := connectionactor.NewConnectionActor("connection", arr)

	if err := dynamicSupervisor.CreateActorPoll(5); err != nil {
		return err
	}

	//my router workeractor
	_ = routeractor.NewRouterActor("router")

	_ = autoscaleractor.NewAutoscalingActor("autoscaling")

	//log.Println(connectionMaker,routerActor,autoscalingActor)

	connectionMaker.ActorLoop()

	return nil
}
