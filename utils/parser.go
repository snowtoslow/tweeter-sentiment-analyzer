package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
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

func AnalyzeSentiments(text string) (counter int8) {
	for _, v := range strings.Fields(text) {
		if val, ok := sentiments.SentimentStorage[v]; ok {
			counter += val
		}
	}
	return
}

func AnalyzeSentimentsTest(text string) (result sentiments.StorageOfSentiments) {
	result = make(map[string]int8)
	var counter int8
	for _, v := range strings.Fields(text) {
		if val, ok := sentiments.SentimentStorage[v]; ok {
			result[v] = val
			counter += val
		}
	}

	result["COUNTER"] = counter
	return
}

func EngagementRatio(retweetedStatus models.RetweetedStatus, favorites, followers int64) (engagementRatio float64) {
	//if retweeted status is nil assign 0;
	//if number of followers is zero return automatically 1
	if followers != 0 {
		engagementRatio = float64((favorites + handleRetweetedStatus(retweetedStatus)) / followers)
	} else {
		engagementRatio = 1
	}
	return
}

func handleRetweetedStatus(retweetedStatus models.RetweetedStatus) (convertedToNr int64) {
	if &retweetedStatus != nil {
		convertedToNr = 1
	}
	return
}

func GetActorPollNameByActorIdentity(identity string) (actorPollName string) {
	if strings.Contains(identity, constants.SentimentActorPool) {
		actorPollName = constants.SentimentActorPool
	} else if strings.Contains(identity, constants.AggregationActorPool) {
		actorPollName = constants.AggregationActorPool
	}
	return
}

func GenerateUuidgen() string {
	f, err := os.Open(constants.NameOfFileToGenerateRandomUuid)
	if err != nil {
		log.Fatal("ERROR OPENING FILE " + constants.NameOfFileToGenerateRandomUuid)
	}
	b := make([]byte, 16)
	if _, err = f.Read(b); err != nil {
		log.Fatal("ERROR READING FROM FILE: " + constants.NameOfFileToGenerateRandomUuid)
	}
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	fmt.Println(uuid)
	return uuid
}

//If someone somewhere at a specific moment of time in our magicGalactic will try to run myMagicProgram and previous method wouldn't work,
//Comment lines 84 - 96 and uncomment from 101 - 115;
/*func GenerateUuidgen() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	log.Println(uuid)


	return
}*/
