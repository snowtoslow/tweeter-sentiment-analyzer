package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"tweeter-sentiment-analyzer/models"
)

func GetRoutes(address string) (mainRoutes *models.MainRouteMsg, err error) {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &mainRoutes)
	if err != nil {
		return
	}

	return
}
