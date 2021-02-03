package main

import (
	"log"
	"runtime"
	"tweeter-sentiment-analyzer/actor-model/actor"
	"tweeter-sentiment-analyzer/utils"
)

func main() {
	log.Println("entry point!")

	mainRouterStruct, err := utils.GetRoutes("http://localhost:4000")
	if err != nil {
		log.Println("ERR OCCURED:", err)
	}

	//director := supervisor.NewSupervisor()

	actorPool := actor.CreateActorPool(5) // actor pool created here!

	runtime.GOMAXPROCS(7)

	chToRecvData := make(chan chan string, 10)

	for _, v := range mainRouterStruct.Routes[:2] {
		go utils.MakeRequest("http://localhost:4000"+v, actorPool, chToRecvData)
	}

	for range chToRecvData {
		_ = <-chToRecvData
	}

}
