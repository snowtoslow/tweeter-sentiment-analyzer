package main

import (
	"log"
	"runtime"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/utils"
)

func main() {
	log.Println("entry point!")

	mainRouterStruct, err := utils.GetRoutes(constants.EndPointToTrigger)
	if err != nil {
		log.Println("ERR OCCURED:", err)
	}

	//director := supervisor.NewSupervisor()
	cpuNumber := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNumber)

	chToRecvData := make(chan string, constants.GlobalChanSiz)

	for _, v := range mainRouterStruct.Routes[:2] {
		go utils.MakeRequest(constants.EndPointToTrigger+v, chToRecvData)
	}

	for range mainRouterStruct.Routes[:2] {
		_ = <-chToRecvData
	}

}
