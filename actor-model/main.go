package main

import (
	"log"
	"runtime"
	"tweeter-sentiment-analyzer/actor-model/actorsystem"
	"tweeter-sentiment-analyzer/actor-model/constants"
	"tweeter-sentiment-analyzer/actor-model/utils"
)

func main() {
	log.Println("entry point!")

	mainRouterStruct, err := utils.GetRoutes(constants.EndPointToTrigger)
	if err == nil {
		log.Println("ERR OCCURRED:", err)
	}

	cpuNumber := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNumber)
	if err = actorsystem.RunApp(mainRouterStruct.Routes[:2]); err != nil {
		log.Println("error occurred after running app:", err)
	}

}
