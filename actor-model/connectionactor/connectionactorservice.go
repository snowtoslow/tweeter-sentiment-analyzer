package connectionactor

import (
	"bufio"
	"log"
	"net/http"
	"tweeter-sentiment-analyzer/constants"
)

func NewConnectionActor(actorName string) *ConnectionActor {
	chanToSendData := make(chan string, constants.GlobalChanSize)

	conenctionMaker := &ConnectionActor{
		Identity:       actorName + constants.ActorName,
		ChanToSendData: chanToSendData,
	}

	// go conenctionMaker.actorLoop()

	return conenctionMaker
}

func (connectionMaker *ConnectionActor) MakeRequest(url string) {
	res, err := http.Get(url)
	if err != nil {
		return
	}

	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		connectionMaker.sendMessage(scanner.Text())
	}
}

func (connectionMaker *ConnectionActor) sendMessage(data string) {
	connectionMaker.ChanToSendData <- data
	log.Println(<-connectionMaker.ChanToSendData)
}

func (connectionMaker *ConnectionActor) ActorLoop(ch chan string) {
	defer close(connectionMaker.ChanToSendData)
	for {
		log.Println(<-connectionMaker.ChanToSendData)
	}
}
