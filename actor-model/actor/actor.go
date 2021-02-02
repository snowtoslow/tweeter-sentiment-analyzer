package actor

type Action func()

type Actor struct {
	Identity   string
	ActionChan chan string
}
