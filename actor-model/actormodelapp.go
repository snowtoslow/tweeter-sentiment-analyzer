package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/constants"
)

func RunApp(arr []string) {

	connectionMaker := connectionactor.NewConnectionActor("connection")

	/*routerActor, err := routeractor.NewRouterActor("router", 5) // here is created router actor which is also a siple actor but which can route messages to actors from pool!
	if err != nil {
		log.Println(err)
	}*/

	for _, v := range arr[:2] {
		go connectionMaker.MakeRequest(constants.EndPointToTrigger + v)
	}

	//connectionMaker.ActorLoop(routerActor.ChanToRecvMsg)

}
