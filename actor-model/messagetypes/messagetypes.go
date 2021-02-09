package messagetypes

type PanicMessage string

type ErrorForSupervisor struct {
	FailedActorIdentity       string
	PanicWithRecoveryFunction func()
}
