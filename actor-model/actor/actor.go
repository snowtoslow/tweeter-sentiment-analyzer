package actor

type Action func()

type Actor struct {
	Address    string
	Identity   string
	ActionChan chan<- string
}
