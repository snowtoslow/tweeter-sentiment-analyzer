package actorregistry

import (
	"fmt"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type ActorRegistry map[string]interface{}

var MyActorRegistry = &ActorRegistry{}

func (registry ActorRegistry) TestFindActorByName(name string) interface{} {
	if x, found := (*MyActorRegistry)[name]; found {
		if res, ok := x.(actorabstraction.IActor); ok {
			return res
		}
		if res, ok := x.([]actorabstraction.IActor); ok {
			return res
		}
	}
	return fmt.Errorf("actor type not found")
}
