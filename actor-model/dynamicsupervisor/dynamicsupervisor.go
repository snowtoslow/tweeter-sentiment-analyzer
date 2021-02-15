package dynamicsupervisor

type DynamicSupervisor struct {
	Identity                            string
	ChanToReceiveNumberOfActorsToCreate chan int
	ChanToReceiveErrors                 chan string
}
