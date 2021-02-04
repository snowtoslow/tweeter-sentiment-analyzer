package actor

type Action func()

type Actor struct {
	Identity          string
	ChanToReceiveData chan interface{}
}
