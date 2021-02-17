package actorabstraction

type IActor interface {
	ActorLoop()
	SendMessage(msg interface{})
}

type AbstractActor struct {
	Identity          string
	ChanToReceiveData chan interface{}
}
