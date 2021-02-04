package actor

import (
	"fmt"
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
	return actorPoll[randomIndex]
}

/*func (actor *Actor) SendProcessedMessage(data string, actors []*Actor) {
	randomActor := GetRandomActor(actors) //pick a random actor from my pool of actors;
	randomActor.ChanToReceiveData <- data        // send msg to this random actor from router actor;
	log.Printf("id:%s---->%s", randomActor.Identity, data)
}*/

func NewActor(actorName string) *Actor {
	chanToRecv := make(chan interface{}, 10)
	actor := &Actor{
		Identity:          actorName + "_actor",
		ChanToReceiveData: chanToRecv,
	}

	go actor.actorLoop(chanToRecv)
	return actor
}

func (actor *Actor) actorLoop(actionChan <-chan interface{}) {
	defer close(actor.ChanToReceiveData)
	for {
		action := <-actor.ChanToReceiveData
		if fmt.Sprintf("%T", action) == "*models.MyJsonName" {
			log.Println("Stuff to count:")
		} else if fmt.Sprintf("%T", action) == "message_types.PanicMessage" {
			log.Println("ERROR:")
		}
		//log.Println("HERE:", action)
	}
}
