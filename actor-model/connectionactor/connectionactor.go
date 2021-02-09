package connectionactor

type ConnectionActor struct {
	Identity       string
	ChanToSendData chan string
}
