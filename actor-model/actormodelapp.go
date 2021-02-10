package actor_model

import (
	"tweeter-sentiment-analyzer/actor-model/actor"
	"tweeter-sentiment-analyzer/actor-model/autoscalleractor"
	"tweeter-sentiment-analyzer/actor-model/connectionactor"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
)

func RunApp(arr []string) error {
	//my connection actor
	connectionMaker := connectionactor.NewConnectionActor("connection", arr)

	//my actor poll
	actorPool, err := actor.CreateActorPoll(5) // actor pool created here!
	if err != nil {
		return err
	}

	//my router actor
	routerActor := routeractor.NewRouterActor("router", actorPool) // here is created router actor which is also a siple actor but which can route messages to actors from pool!

	//connectionMaker.SendDataToDifferentActorsOverChan(routerActor.ChanToRecvMsg) //before commenting SendDatToChan

	autoscallerActor := autoscalleractor.NewAutoscallerActor("autoscaller")

	connectionMaker.SendDataToMultipleActorsOverChan(routerActor.ChanToRecvMsg, autoscallerActor.ChanToReceiveMessagesForCount)

	return nil
}

//stuff before add merge chan method;
//connectionMaker.ReceivePreparedData(arr[:2], routerActor.ChanToRecvMsg)
