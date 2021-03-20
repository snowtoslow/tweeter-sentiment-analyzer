package main

import (
	"log"
	"runtime"
	actor_model "tweeter-sentiment-analyzer/actor-model"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func main() {
	log.Println("entry point!")

	mainRouterStruct, err := utils.GetRoutes(constants.EndPointToTrigger)
	if err != nil {
		log.Println("ERR OCCURED:", err)
	}

	cpuNumber := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNumber)
	if err = actor_model.RunApp(mainRouterStruct.Routes[:2]); err != nil {
		log.Println("error occurred after running app:", err)
	}

}
