package connectionactor

import (
	"bufio"
	"net/http"
	"regexp"
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

func (connectionMaker *ConnectionActor) SendDataToMultipleActorsOverChan(cs ...chan string) {
	for msg := range connectionMaker.receivePreparedData(connectionMaker.AddressRoutesArray) {
		for range cs {
			select {
			case cs[0] <- msg:
				cs[0] <- connectionMaker.createMessage(msg)
			case cs[1] <- msg:
				cs[1] <- regexp.MustCompile(constants.JsonRegex).FindString(msg)
			}
			//v <- connectionMaker.createMessage(msg)
		}
	}
}

func (connectionMaker *ConnectionActor) SendDataToActorChan(ch chan string) {
	for msg := range connectionMaker.receivePreparedData(connectionMaker.AddressRoutesArray) {
		ch <- connectionMaker.createMessage(msg)
	}
}

func (connectionMaker *ConnectionActor) receivePreparedData(arr []string) chan string {
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
