package routerstrategy

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type RoundRobin struct {
	ActorsPool   []actorabstraction.IActor
	CurrentIndex int
}

func NewRoundRobinStrategy(actors []actorabstraction.IActor) *RoundRobin {
	return &RoundRobin{
		ActorsPool:   actors,
		CurrentIndex: 0,
	}
}

func MultipleBalancerEntity(arrayOfArray ...[]actorabstraction.IActor) []*RoundRobin {
	var balancers []*RoundRobin
	for _, v := range arrayOfArray {
		balancers = append(balancers, NewRoundRobinStrategy(v))
	}
	return balancers
}

func (r *RoundRobin) Balancer(msg interface{}) {
	if r.CurrentIndex >= len(r.ActorsPool) {
		r.CurrentIndex = 0
	}
	r.ActorsPool[r.CurrentIndex].SendMessage(msg)
	//log.Printf("%+v\n",r.ActorsPool[r.CurrentIndex])
	r.CurrentIndex++
}
