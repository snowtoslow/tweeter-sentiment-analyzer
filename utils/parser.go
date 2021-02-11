package utils

import (
	"encoding/json"
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
