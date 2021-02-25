package routerstrategy

import (
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type RoundRobin struct {
	ActorsPool   []actorabstraction.IActor
	CurrentIndex int
}

func NewRoundRobinStrategy() *RoundRobin {
	return &RoundRobin{
		CurrentIndex: 0,
	}
}

func (r *RoundRobin) MultipleBalancerEntity(arrayOfArray ...[]actorabstraction.IActor) []*RoundRobin {
	var balancers []*RoundRobin
	for _, v := range arrayOfArray {
		balancer := NewRoundRobinStrategy()
		balancer.ActorsPool = v
		balancers = append(balancers, balancer)
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
