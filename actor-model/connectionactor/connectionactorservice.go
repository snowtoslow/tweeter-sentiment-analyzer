package connectionactor

import (
	"bufio"
	"log"
	"net/http"
	"regexp"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewConnectionActor(actorName string) *ConnectionActor {
	chanToSendData := make(chan string, constants.GlobalChanSize)

	connectionMaker := &ConnectionActor{
		Identity:       actorName + constants.ActorName,
		ChanToSendData: chanToSendData,
	}

	(*actorregistry.MyActorRegistry)["connectionActor"] = connectionMaker

	//go connectionMaker.ActorLoop()
	//we can uncomment it but the could ud become more harder wo change because we need the ability to send to a new chan

	return connectionMaker
}

func (connectionMaker *ConnectionActor) SendDataToMultipleActorsOverChan(routes []string, cs ...chan string) {
	for msg := range connectionMaker.receivePreparedData(routes) {
		for _, v := range cs {
			v <- connectionMaker.createMessage(msg)
		}
	}
}

func (connectionMaker *ConnectionActor) receivePreparedData(arr []string) chan string {
	arrayOfChannels := make([]chan string, len(arr))
	for k, v := range arr {
		arrayOfChannels[k] = connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v))
	}
	return utils.MergeWait(arrayOfChannels...)
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

func (connectionMaker *ConnectionActor) createMessage(data string) string {
	connectionMaker.ChanToSendData <- data
	return <-connectionMaker.ChanToSendData
}

func (connectionMaker *ConnectionActor) SendMessage(data string) {
	connectionMaker.ChanToSendData <- data
}

func (connectionMaker *ConnectionActor) ActorLoop() {
	defer close(connectionMaker.ChanToSendData)
	regexData := regexp.MustCompile(constants.JsonRegex)
	for {
		receivedString := regexData.FindString(<-connectionMaker.ChanToSendData)
		if receivedString == constants.PanicMessage {
			log.Println("ERROR!")
		}
	}
}
