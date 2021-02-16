package routeractor

type RouterActor struct {
	Identity          string
	ChanToRecvMsg     chan string
	CurrentActorIndex int
}
