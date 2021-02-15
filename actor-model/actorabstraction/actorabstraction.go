package actorabstraction

import "log"

type Actor struct {
	Identity       string
	ChanToSendData chan string
}

func (actor *Actor) actorLoop(actorName string) {
	log.Println("actor loop", actorName)
}
