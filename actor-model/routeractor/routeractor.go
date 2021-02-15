package routeractor

import "tweeter-sentiment-analyzer/actor-model/workeractor"

type RouterActor struct {
	Identity          string
	ChanToRecvMsg     chan string
	CurrentActorIndex int
	Actors            *[]workeractor.Actor
}
