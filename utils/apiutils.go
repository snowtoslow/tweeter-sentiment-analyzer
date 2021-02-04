package utils

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"tweeter-sentiment-analyzer/actor-model/routeractor"
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

func MakeRequest(url string, ch chan string) {
	res, err := http.Get(url)
	if err != nil {
		close(ch)
		return
	}
	routerActor, err := routeractor.NewRouterActor("router", 5) // here is created router actor which is also a siple actor but which can route messages to actors from pool!
	if err != nil {
		close(ch)
		return
	}

	defer close(ch)
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		routerActor.SendMessage(scanner.Text())
	}
}
