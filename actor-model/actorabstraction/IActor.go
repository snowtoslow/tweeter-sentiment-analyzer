package actorabstraction

type IActor interface {
	ActorLoop()
	SendMessage(msg interface{})
}

type AbstractActor struct {
	IActor
	Identity          string
	ChanToReceiveData chan string
}
