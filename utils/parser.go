package utils

import (
	"encoding/json"
	"math"
	msgType "tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
)

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
