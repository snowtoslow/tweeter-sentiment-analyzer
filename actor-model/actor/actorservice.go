package actor

func (actor *Actor) SendMessage() (err error) {
	return err
}

func (actor *Actor) ActorLoop() (err error) {
	return err
}

func (actor *Actor) NewActor() *Actor {
	return &Actor{
		Address:  "actor1",
		Identity: "identity1",
	}
}
