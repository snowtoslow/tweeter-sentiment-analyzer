package main

import (
	"log"
	"runtime"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
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

	chToRecvData := make(chan string, constants.GlobalChanSize)

	routerActor, err := routeractor.NewRouterActor("router", 5) // here is created router actor which is also a siple actor but which can route messages to actors from pool!
	if err != nil {
		log.Println(err)
	}

	for _, v := range mainRouterStruct.Routes[:2] {
		go routerActor.MakeRequest(constants.EndPointToTrigger+v, chToRecvData)
	}

	for {
		log.Println(<-chToRecvData)
	}

}

/*connectiin = MakeConnActor
router = MakeRouterActor
DynSup = Make
AutoScale = Make

Reg = Make
Reg.register(...)
AppSup.link*/
