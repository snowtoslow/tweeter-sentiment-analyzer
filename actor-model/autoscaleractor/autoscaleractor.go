package autoscaleractor

type AutoscalingActor struct {
	Identity                      string
	ChanToReceiveMessagesForCount chan string
	ChanToSendCounterResult       chan int
}
