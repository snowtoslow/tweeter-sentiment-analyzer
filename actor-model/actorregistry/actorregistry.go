package actorregistry

import (
	"fmt"
)

type ActorRegistry map[string]interface{}

var MyActorRegistry = &ActorRegistry{}

func (registry ActorRegistry) TestFindActorByName(name string) interface{} {
	if x, found := (*MyActorRegistry)[name]; found {
		/*if res, ok := x.(*actorabstraction.IActor); ok {
			log.Println("actor:", res)
			return res
		}
		if res, ok := x.([]actorabstraction.IActor); ok {
			log.Println("actors pool:", res)
			return res
		}*/
		//fmt.Printf("%T",x)
		return x
	}
	return fmt.Errorf("actor type not found")
}
