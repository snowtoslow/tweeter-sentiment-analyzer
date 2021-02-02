package actor

type Action func()

type Actor struct {
	// Address    string
	IsBusy     bool
	Identity   string
	ActionChan chan<- string
}
