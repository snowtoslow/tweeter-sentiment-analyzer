package supervisor

type Supervisor struct {
	ActorAddresses *[]string
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		ActorAddresses: new([]string),
	}
}
