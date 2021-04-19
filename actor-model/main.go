package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"tweeter-sentiment-analyzer/actor-model/actorsystem"
	"tweeter-sentiment-analyzer/actor-model/utils"
)

func main() {
	log.Println("entry point!")
	//dev purpose
	os.Setenv("RTP_SERVER", "localhost:4000")
	os.Setenv("BROKER_URL", "localhost:8088")

	rtpImageUl := fmt.Sprintf("http://%s", os.Getenv("RTP_SERVER")) //change here from constants.EndPointToTrigger

	mainRouterStruct, err := utils.GetRoutes(rtpImageUl) //change here!
	if err != nil {
		log.Fatalf("Error occured getting addresses of routes: %s", err)
		return
	}

	cpuNumber := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNumber)
	if err = actorsystem.RunApp(mainRouterStruct.Routes[:2]); err != nil {
		log.Fatalf("Error occured running applciation: %s", err)
		return
	}

}
