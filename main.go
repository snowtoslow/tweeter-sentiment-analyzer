package main

import (
	"log"
	"runtime"
	"tweeter-sentiment-analyzer/actor-model/actor"
	"tweeter-sentiment-analyzer/utils"
)

func main() {
	log.Println("entry point!")

	_, err := utils.GetRoutes("http://localhost:4000")
	if err != nil {
		log.Println("ERR OCCURED:", err)
	}

	//director := supervisor.NewSupervisor()

	actorPool := actor.CreateActorPool(5) // actor pool created here!

	runtime.GOMAXPROCS(7)

	utils.MakeRequest("http://localhost:4000/tweets/1", actorPool)

}
