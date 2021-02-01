package actor

import (
	"log"
	"strconv"
)

func CreateActorPool(numberOfActors int) (actorPoll []*Actor) {
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, NewActor(i))
	}
	return actorPoll
}

func (actor *Actor) SendMessage(data string) {
	actor.ActionChan <- data
}

func NewActor(actorNum int) *Actor {
	chanToRecv := make(chan string, 10)
	actor := &Actor{
		Address:    "actor_" + strconv.Itoa(actorNum),
		Identity:   "ident_" + strconv.Itoa(actorNum),
		ActionChan: chanToRecv,
	}

	go actor.actorLoop(chanToRecv)
	return actor
}

func (actor *Actor) actorLoop(actionChan <-chan string) {
	defer close(actor.ActionChan)
	for {
		action := <-actionChan
		log.Println(action)
	}
}
