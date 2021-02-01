package actor

import (
	"log"
	"strconv"
)

/*func CreateActorPool(numberOfActors int)(actorPoll []*Actor){
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, NewActor(i))
	}
	return actorPoll
}*/

func (actor *Actor) SendMessage() (err error) {
	return err
}

func NewActor(actorNum int, chanToRecv chan string) *Actor {
	actor := &Actor{
		Address:    "actor_" + strconv.Itoa(actorNum),
		Identity:   "ident_" + strconv.Itoa(actorNum),
		ActionChan: chanToRecv,
	}

	go actor.actorLoop(chanToRecv)
	return actor
}

func (actor *Actor) actorLoop(actionChan chan string) {
	defer close(actionChan)
	for {
		action := <-actionChan
		log.Println("PRINT MY ACTION:", action)
	}
}
