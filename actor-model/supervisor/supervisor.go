package supervisor

import (
	"tweeter-sentiment-analyzer/actor-model/messagetypes"
)

type Supervisor struct {
	Identity                  string
	ChanToReceiveErrorMessage chan *messagetypes.ErrorForSupervisor
}
