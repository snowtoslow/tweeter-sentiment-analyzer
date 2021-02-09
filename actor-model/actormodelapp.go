package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
)

func RunApp(arr []string) error {

	connectionMaker := connectionactor.NewConnectionActor("connection")

	routerActor, err := routeractor.NewRouterActor("router", 5) // here is created router actor which is also a siple actor but which can route messages to actors from pool!
	if err != nil {
		return err
	}

	//connectionMaker.SendDataToConnectionActor(connectionMaker.ReceivePreparedData(arr[:2])) before commenting SendDatToChan
	connectionMaker.ReceivePreparedData(arr[:2], routerActor.ChanToRecvMsg)

	return nil
}
