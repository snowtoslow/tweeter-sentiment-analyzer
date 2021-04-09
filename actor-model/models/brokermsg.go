package models

import (
	"tweeter-sentiment-analyzer/actor-model/constants"
)

type UserTopic string

type TweetsTopic string

type BroadCastTopic string

type BrokerMsg struct {
	Content interface{} `json:"content,omitempty"` //  tweets or users;
	Topic   interface{} `json:"topic,omitempty"`   // one of thee topic types
}

func (bk *BrokerMsg) SetTopic(actionType string) {
	if actionType == constants.JsonNameOfStruct || actionType == constants.RetweetedStatus {
		bk.Topic = TweetsTopic(constants.TweetsTopic)
	} else {
		bk.Topic = UserTopic(constants.UserTopic)
	}
}
