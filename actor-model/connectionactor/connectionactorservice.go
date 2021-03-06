package connectionactor

import (
	"bufio"
	"net/http"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/actor-model/autoscaleractor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func NewConnectionActor(actorName string, routes []string) actorabstraction.IActor {
	chanToSendData := make(chan interface{}, constants.GlobalChanSize)

	connectionMaker := &ConnectionActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToSendData,
		},
		Routes: routes,
	}

	(*actorregistry.MyActorRegistry)["connectionActor"] = connectionMaker

	//go connectionMaker.ActorProps.ActorLoop()
	//we can uncomment it but the could ud become more harder wo change because we need the ability to send to a new chan

	return connectionMaker
}

func (connectionMaker *ConnectionActor) ActorLoop() {
	cs := []actorabstraction.IActor{
		actorregistry.MyActorRegistry.FindActorByName("routerActor").(*routeractor.RouterActor),
		actorregistry.MyActorRegistry.FindActorByName("autoscalingActor").(*autoscaleractor.AutoscalingActor),
	}
	for msg := range connectionMaker.receivePreparedData(connectionMaker.Routes) {
		for _, v := range cs {
			if len(msg) != 0 {
				v.SendMessage(msg)
			}
			/*connectionMaker.SendMessage(msg)
			v <- connectionMaker.getChan()*/
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

func (connectionMaker *ConnectionActor) SendMessage(data interface{}) {
	connectionMaker.ActorProps.ChanToReceiveData <- data
}

/*func (connectionMaker *ConnectionActor) ActorLoop() {
	defer close(connectionMaker.ActorProps.ChanToReceiveData)
	regexData := regexp.MustCompile(constants.JsonRegex)
	for {
		receivedString := regexData.FindString(<-connectionMaker.ActorProps.ChanToReceiveData)
		if receivedString == constants.PanicMessage {
			log.Println("ERROR!")
		}
		log.Println(<-connectionMaker.ActorProps.ChanToReceiveData)
	}
}*/
