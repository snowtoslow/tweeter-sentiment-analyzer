package connectionactor

import (
	"bufio"
	"log"
	"net/http"
	"regexp"
	"tweeter-sentiment-analyzer/constants"
)

func NewConnectionActor(actorName string) *ConnectionActor {
	chanToSendData := make(chan string, constants.GlobalChanSize)

	connectionMaker := &ConnectionActor{
		Identity:       actorName + constants.ActorName,
		ChanToSendData: chanToSendData,
	}

	//go connectionMaker.actorLoop()
	//we can uncomment it but the could ud become more harder wo change because we need the ability to send to a new chan

	return connectionMaker
}

/*func (connectionMaker *ConnectionActor) SendDataToConnectionActor(recch <-chan string) {
	for msg := range recch {
		connectionMaker.sendMessage(msg)
	}
}*/

func (connectionMaker *ConnectionActor) ReceivePreparedData(arr []string, ch chan string) {
	//c := make(chan string, 10)
	for _, v := range arr[:2] {
		/*c = connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v))*/
		for msg := range connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v)) {
			ch <- connectionMaker.createMessage(msg)
		}
	}

	/*for msg := range c {
		connectionMaker.sendMessage(msg)
	}*/
	//return c
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

func (connectionMaker *ConnectionActor) sendMessage(data string) {
	connectionMaker.ChanToSendData <- data
}

func (connectionMaker *ConnectionActor) actorLoop() {
	defer close(connectionMaker.ChanToSendData)
	regexData := regexp.MustCompile(constants.JsonRegex)
	for {
		receivedString := regexData.FindString(<-connectionMaker.ChanToSendData)
		if receivedString == constants.PanicMessage {
			log.Println("ERROR!")
		}
	}
}
