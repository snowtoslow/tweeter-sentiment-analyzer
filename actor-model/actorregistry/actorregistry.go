package actorregistry

import (
	"fmt"
	"log"
	"tweeter-sentiment-analyzer/actor-model/autoscaleractor"
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/actor-model/dynamicsupervisor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
	worker_actor "tweeter-sentiment-analyzer/actor-model/workeractor"
)

type ActorRegistry map[string]interface{}

var MyActorRegistry = ActorRegistry{
	"dynamicSupervisor": dynamicsupervisor.NewDynamicSupervisor("dynamic_supervisor"),
	"connectionMaker":   connectionactor.NewConnectionActor("connection"),
	/*"autoscalingActor" : autoscaleractor.NewAutoscalingActor("autoscaling"),
	"routerActor" : routeractor.NewRouterActor("router"),*/
	/*"actorPoll" : dynamicsupervisor.NewDynamicSupervisor("dynamic_supervisor").CreateActorPoll(5),*/
}

func (registry ActorRegistry) TestFindActorByName(name string) interface{} {
	if x, found := MyActorRegistry[name]; found {
		if res, ok := x.(*dynamicsupervisor.DynamicSupervisor); ok {
			log.Println("Dynamic supervisor:", res)
			return res
		}
		if res, ok := x.(*connectionactor.ConnectionActor); ok {
			log.Println("connection maker:", res)
			return res
		}
		if res, ok := x.(*autoscaleractor.AutoscalingActor); !ok {
			log.Println("autoscaling actor:", res)
			return res
		}
		if res, ok := x.(*routeractor.RouterActor); !ok {
			log.Println("router actor", res)
			return res
		}
		if res, ok := x.(*[]worker_actor.Actor); !ok {
			log.Println("worker actors:", res)
			return res
		}
	}
	return fmt.Errorf("actor type not found")
}
