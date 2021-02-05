package actor

import "tweeter-sentiment-analyzer/actor-model/messagetypes"

type Actor struct {
	Identity          string
	ChanToReceiveData chan string
	ChanToSendError   chan *messagetypes.ErrorForSupervisor
}
