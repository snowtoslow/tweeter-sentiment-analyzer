package connectionactor

import (
	"bufio"
	"log"
	"net/http"
	"tweeter-sentiment-analyzer/constants"
)

func NewConnectionActor(actorName string) *ConnectionActor {
	chanToSendData := make(chan string, constants.GlobalChanSize)

	connectionMaker := &ConnectionActor{
		Identity:       actorName + constants.ActorName,
		ChanToSendData: chanToSendData,
	}

	go connectionMaker.actorLoop()

	return connectionMaker
}

func (connectionMaker *ConnectionActor) SendDataToConnectionActor(recch <-chan string) {
	for msg := range recch {
		connectionMaker.sendMessage(msg)
	}
}

func (connectionMaker *ConnectionActor) ReceivePreparedData(arr []string) chan string {
	c := make(chan string, 10)
	for _, v := range arr[:2] {
		c = connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v))
	}
	return c
}

func (connectionMaker *ConnectionActor) getPreparedData(ic <-chan string) chan string {
	oc := make(chan string, constants.GlobalChanSize)
	go func() {
		for v := range ic {
			oc <- v
		}
		close(oc)
	}()
	return oc
}

func (connectionMaker *ConnectionActor) makeReqPipeline(url string) chan string {
	dataFlowChan := make(chan string, constants.GlobalChanSize)
	go func() {
		res, err := http.Get(constants.EndPointToTrigger + url)
		if err != nil {
			return
		}
		defer res.Body.Close()
		scanner := bufio.NewScanner(res.Body)

		for scanner.Scan() {
			dataFlowChan <- scanner.Text()
		}
		close(dataFlowChan)
	}()
	return dataFlowChan
}

func (connectionMaker *ConnectionActor) sendMessage(data string) {
	connectionMaker.ChanToSendData <- data
}

func (connectionMaker *ConnectionActor) actorLoop() {
	defer close(connectionMaker.ChanToSendData)
	for {
		log.Println("HERE:", <-connectionMaker.ChanToSendData)
	}
}
