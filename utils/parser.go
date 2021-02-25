package utils

import (
	"encoding/json"
	"math"
	"strings"
	"tweeter-sentiment-analyzer/actor-model/message-types"
	"tweeter-sentiment-analyzer/constants"
	"tweeter-sentiment-analyzer/models"
	"tweeter-sentiment-analyzer/sentiments"
)

func CreateMessageType(processedString string) interface{} {
	if processedString == constants.PanicMessage {
		return message_types.PanicMessage(processedString)
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

func AnalyzeSentiments(text string) (result sentiments.StorageOfSentiments) {
	result = make(map[string]int8)
	var counter int8
	for _, v := range strings.Fields(text) {
		if val, ok := sentiments.SentimentStorage[v]; ok {
			result[v] = val
			counter += val
		}
	}

	result[constants.CounterFinalResult] = counter
	return
}

func ComputeEngagementRatio(favorites, retweets, followers int) float32 {
	return 0
}
