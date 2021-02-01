package utils

import (
	"encoding/json"
	"io"
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

func GetStreams(address string, ch chan string) (err error) {
	res, err := http.Get(address)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	p := make([]byte, 4)
	for {
		_, err := res.Body.Read(p)
		if err == io.EOF {
			break
		}
	}

	log.Fatal(string(p))

	return nil
}

func MakeRequest(url string, ch chan string) {
	res, err := http.Get(url)
	if err != nil {
		close(ch)
		return
	}
	data := make([]byte, 128)
	defer res.Body.Close()
	defer close(ch)

	for n, err := res.Body.Read(data); err == nil; n, err = res.Body.Read(data) {
		ch <- string(data[:n])
		actor.NewActor(1, ch)
	}
}
