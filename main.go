package main

import (
	"log"
	"runtime"
	"tweeter-sentiment-analyzer/utils"
)

func main() {
	log.Println("entry point!")
	/*
		_, err := utils.GetRoutes("http://localhost:4000")
		if err!=nil {
			log.Println("ERR OCCURED:",err)
		}*/

	//director := supervisor.NewSupervisor()

	/*for _, v := range routesStruct.Routes{

	}*/

	runtime.GOMAXPROCS(7)

	ch1 := make(chan string, 128)

	utils.MakeRequest("http://localhost:4000/tweets/1", ch1)

}
