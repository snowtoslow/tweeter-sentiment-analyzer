package autoscalleractor

type AutoscallerActor struct {
	Identity                      string
	ChanToReceiveMessagesForCount chan string
	ChanToSendCounterResult       chan int
}
