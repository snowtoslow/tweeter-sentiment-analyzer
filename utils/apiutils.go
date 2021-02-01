package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"tweeter-sentiment-analyzer/actor-model/actor"
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
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &mainRoutes)
	if err != nil {
		return nil, err
	}

	return mainRoutes, nil
}

func MakeRequest(url string, ch chan string) {
	res, err := http.Get(url)
	if err != nil {
		close(ch)
		return
	}
	data := make([]byte, 512)
	defer res.Body.Close()

	for n, err := res.Body.Read(data); err == nil; n, err = res.Body.Read(data) {
		myActor := actor.NewActor(1)
		myActor.SendMessage(string(data[:n]))
	}
}
