package message_types

type PanicMessage string

type ErrorToSupervisor struct {
	ActorIdentity string
	Message       PanicMessage
}
