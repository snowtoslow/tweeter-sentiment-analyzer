package workeractor

type Actor struct {
	Identity          string
	ChanToReceiveData chan string
}
