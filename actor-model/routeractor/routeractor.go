package routeractor

import "tweeter-sentiment-analyzer/actor-model/actor"

type RouterActor struct {
	Identity          string
	ChanToRecvMsg     chan string
	CurrentActorIndex int
	Actors            *[]actor.Actor
}
