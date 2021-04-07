package clientactor

import (
	"net"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
)

type ClientActor struct {
	ActorProps actorabstraction.AbstractActor
	Connection net.Dialer
}
