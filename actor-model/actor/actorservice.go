package actor

import (
	"log"
	"math/rand"
	"strconv"
)

func CreateActorPool(numberOfActors int) (actorPoll []*Actor) {
	for i := 0; i < numberOfActors; i++ {
		actorPoll = append(actorPoll, NewActor("working_"+strconv.Itoa(i)))
	}
	return actorPoll
}

//generate random actors from my array
func GetRandomActor(actorPoll []*Actor) *Actor {
	randomIndex := rand.Intn(len(actorPoll))
	log.Println("Actor index:", randomIndex)
	return actorPoll[randomIndex]
}

/*func (actor *Actor) SendMessage(data string) {
	actor.ActionChan <- data
}*/

func (actor *Actor) SendMessage(data string, actors []*Actor) {
	randomActor := GetRandomActor(actors) //pick a random actor from my pool of actors;
	randomActor.ActionChan <- data        // send msg to this random actor from router actor;
	log.Printf("id:%s---->%s", randomActor.Identity, data)
}

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
