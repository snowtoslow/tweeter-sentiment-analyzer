package constants

var (
	PanicMessage           = "{\"message\": panic}"
	JsonRegex              = "\\{.*\\:\\{.*\\:.*\\}\\}|\\{(.*?)\\}"
	EndPointToTrigger      = "http://localhost:4000"
	JsonNameOfStruct       = "*models.MyJsonName"
	RetweetedStatus        = "models.RetweetedStatus"
	PanicMessageType       = "message_types.PanicMessage"
	ActorName              = "_actor_"
	GlobalChanSize         = 10
	PointerJsonSyntaxError = "*json.SyntaxError"
	DefaultActorPollSize   = 5
	CounterFinalResult     = "COUNTER_RESULT"
	SentimentActorPool     = "sentimentActorPool"
	AggregationActorPool   = "aggregationActorPool"
)
