package connectionactor

import (
	"bufio"
	"net/http"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewConnectionActor(actorName string, routesArray []string) *ConnectionActor {
	chanToSendData := make(chan string, constants.GlobalChanSize)

	connectionMaker := &ConnectionActor{
		Identity:           actorName + constants.ActorName,
		ChanToSendData:     chanToSendData,
		AddressRoutesArray: routesArray,
	}

	//go connectionMaker.actorLoop()
	//we can uncomment it but the could ud become more harder wo change because we need the ability to send to a new chan

	return connectionMaker
}

func (connectionMaker *ConnectionActor) SendDataToDifferentActorsOverChan(ch chan string) {
	for msg := range connectionMaker.receivePreparedDataTest(connectionMaker.AddressRoutesArray) {
		ch <- connectionMaker.createMessage(msg)
	}
}

func (connectionMaker *ConnectionActor) receivePreparedDataTest(arr []string) chan string {
	chansArr := make([]chan string, len(arr))
	for k, v := range arr {
		chansArr[k] = connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v))
	}
	return utils.MergeWait(chansArr...)
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

/*func (connectionMaker *ConnectionActor) sendMessage(data string) {
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

//received prepared data method before adding merge of two chan function;
func (connectionMaker *ConnectionActor) ReceivePreparedData(arr []string) chan string {
	c := make(chan string, 10)
	for _, v := range arr[:2] {
		c = connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v))
	}

	return c
}*/
