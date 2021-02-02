package actor

import (
	"log"
	"strconv"
)

func CreateActorPool(numberOfActors int) (actorPoll []*Actor) {
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, NewActor("working_"+strconv.Itoa(i)))
	}
	return actorPoll
}

func (actor *Actor) SendMessage(data string) {
	actor.ActionChan <- data
}

/*func (actor *Actor) SendMessage(data string,actorToRecvMsg *Actor) {
	log.Printf("id:%s---->%s",actorToRecvMsg.Identity,data)
	actorToRecvMsg.ActionChan <- data
}*/

func NewActor(actorName string) *Actor {
	chanToRecv := make(chan string, 10)
	actor := &Actor{
		Identity:   actorName + "_actor",
		ActionChan: chanToRecv,
		IsBusy:     false,
	}

	go actor.actorLoop(chanToRecv)
	return actor
}

func (actor *Actor) actorLoop(actionChan <-chan string) {
	defer close(actor.ActionChan)
	for {
		action := <-actionChan
		log.Printf(action)
	}
}
