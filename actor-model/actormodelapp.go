package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/autoscaleractor"
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/actor-model/dynamicsupervisor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
	"tweeter-sentiment-analyzer/actor-model/sinkactor"
)

func RunApp(arr []string) error {

	if _, err := dynamicsupervisor.NewDynamicSupervisor("dynamic_supervisor"); err != nil {
		return err
	}
	//my connection workeractor
	connectionMaker := connectionactor.NewConnectionActor("connection", arr)

	//my router workeractor
	_ = routeractor.NewRouterActor("router")

	_ = autoscaleractor.NewAutoscalingActor("autoscaling")

	//log.Println(connectionMaker,routerActor,autoscalingActor)

	_ = sinkactor.NewSinkActor("sink")

	connectionMaker.ActorLoop()

	return nil
}
