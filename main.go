package main

import (
	"log"
	"runtime"
	actor_model "tweeter-sentiment-analyzer/actor-model"
	"tweeter-sentiment-analyzer/actor-model/constants"
	"tweeter-sentiment-analyzer/actor-model/utils"
)

func main() {
	log.Println("entry point!")

	mainRouterStruct, err := utils.GetRoutes(constants.EndPointToTrigger)
	if err != nil {
		log.Println("ERR OCCURRED:", err)
	}

	cpuNumber := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNumber)
	if err = actor_model.RunApp(mainRouterStruct.Routes[:2]); err != nil {
		log.Println("error occurred after running app:", err)
	}

}
