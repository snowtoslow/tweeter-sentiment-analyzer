package actorabstraction

type IActor interface {
	ActorLoop()
	SendMessage(msg interface{})
}
