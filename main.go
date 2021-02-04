package main

import (
	"log"
	"runtime"
	"tweeter-sentiment-analyzer/utils"
)

func main() {
	log.Println("entry point!")

	mainRouterStruct, err := utils.GetRoutes("http://localhost:4000")
	if err != nil {
		log.Println("ERR OCCURED:", err)
	}

	//director := supervisor.NewSupervisor()

	runtime.GOMAXPROCS(7)

	chToRecvData := make(chan string, 10)

	for _, v := range mainRouterStruct.Routes[:2] {
		go utils.MakeRequest("http://localhost:4000"+v, chToRecvData)
	}

	for range mainRouterStruct.Routes[:2] {
		_ = <-chToRecvData
	}

}
