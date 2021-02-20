package actorregistry

import (
	"fmt"
	"log"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type ActorRegistry map[string]interface{}

var MyActorRegistry = &ActorRegistry{}

func (registry ActorRegistry) FindActorByName(name string) interface{} {
	if x, found := (*MyActorRegistry)[name]; found {
		if res, ok := x.(actorabstraction.IActor); ok {
			return res
		}
		if res, ok := x.(*[]actorabstraction.IActor); ok {
			return res
		}
		log.Printf("%T\n", x)
	}
	return fmt.Errorf("actor type not found")
}
