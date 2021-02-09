package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
)

func RunApp(arr []string) {

	connectionMaker := connectionactor.NewConnectionActor("connection")

	/*routerActor, err := routeractor.NewRouterActor("router", 5) // here is created router actor which is also a siple actor but which can route messages to actors from pool!
	if err != nil {
		log.Println(err)
	}*/

	/*c := make(chan string,10)
	for _,v :=range arr[:2]{
		c = connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v))
	}

	for data := range c {
		log.Printf("Items saved: %+v", data)
	}*/

	/*	for _,v :=range arr[:2]{
		connectionMaker.SendDataToConnectionActor(connectionMaker.getPreparedData(connectionMaker.makeReqPipeline(v)))
	}*/

	connectionMaker.SendDataToConnectionActor(connectionMaker.ReceivePreparedData(arr[:2]))

}
