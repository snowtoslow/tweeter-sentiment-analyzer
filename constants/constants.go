package constants

var (
	PanicMessage      = "{\"message\": panic}"
	JsonRegex         = "\\{.*\\:\\{.*\\:.*\\}\\}|\\{(.*?)\\}"
	EndPointToTrigger = "http://localhost:4000"
	JsonNameOfStruct  = "*models.MyJsonName"
	PanicMessageType  = "messagetypes.PanicMessage"
	ActorName         = "_actor"
	GlobalChanSize    = 10
)
