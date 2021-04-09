package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"tweeter-sentiment-analyzer/actor-model/models"
)

func GetRoutes(address string) (mainRoutes *models.MainRouteMsg, err error) {
	response, err := http.Get(address)
	if err != nil {
		log.Println("first err:", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	err = json.Unmarshal(body, &mainRoutes)
	if err != nil {
		log.Println("third err:", err)
		return
	}

	return
}
