package actor

type Actor struct {
	Identity          string
	ChanToReceiveData chan string
}
