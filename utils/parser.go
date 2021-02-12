package utils

import (
	"encoding/json"
	"math"
	"regexp"
	"strings"
	msgType "tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
)

func GetChanData(actorChan chan string) interface{} {
	regexData := regexp.MustCompile("data: {(.*?)}")
	receivedString := strings.Split(regexData.FindString(<-actorChan), ":")[1]
	var tweetMsg *models.MyJsonName
	if receivedString == constants.PanicMessage {
		return msgType.PanicMessage(strings.Split(receivedString, ":")[1])
	} else {
		err := json.Unmarshal([]byte(receivedString), &tweetMsg)
		if err != nil {
			return err
		}
		return tweetMsg
	}
}

func CreateMessageType(processedString string) interface{} {
	if processedString == constants.PanicMessage {
		return msgType.PanicMessage(processedString)
	} else {
		var tweetMsg *models.MyJsonName
		if err := json.Unmarshal([]byte(processedString), &tweetMsg); err != nil {
			return err
		}
		return tweetMsg
	}
}

func MovingExpAvg(value, oldValue, fdtime, ftime float64) float64 {
	alpha := 1.0 - math.Exp(-fdtime/ftime)
	r := alpha*value + (1.0-alpha)*oldValue
	return r
}
