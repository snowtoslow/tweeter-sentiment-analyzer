package actor

import (
	"log"
	"math/rand"
	"regexp"
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

/*func (actor *Actor) SendMessage(data string, actors []*Actor) {
	randomActor := GetRandomActor(actors) //pick a random actor from my pool of actors;
	randomActor.ChanToReceiveData <- data        // send msg to this random actor from router actor;
	log.Printf("id:%s---->%s", randomActor.Identity, data)
}*/

func NewActor(actorName string) *Actor {
	chanToRecv := make(chan string, 10)
	actor := &Actor{
		Identity:          actorName + "_actor",
		ChanToReceiveData: chanToRecv,
	}

	go actor.actorLoop(chanToRecv)
	return actor
}

func (actor *Actor) actorLoop(actionChan <-chan string) {
	defer close(actor.ChanToReceiveData)
	/// regexData := regexp.MustCompile("data: {(.*?)}")// already tested
	messageRegex := regexp.MustCompile("event:(.*)") // need to test
	for {
		action := <-actor.ChanToReceiveData
		log.Println("HERE:", messageRegex.FindString(action))
	}
}
