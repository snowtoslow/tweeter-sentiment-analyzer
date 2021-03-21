package constants

import "time"

const (
	PanicMessage           = "{\"message\": panic}"
	JsonRegex              = "\\{.*\\:\\{.*\\:.*\\}\\}|\\{(.*?)\\}"
	EndPointToTrigger      = "http://localhost:4000"
	JsonNameOfStruct       = "*models.MyJsonName"
	RetweetedStatus        = "models.RetweetedStatus"
	RetweetedStatusPointer = "*models.RetweetedStatus"
	PanicMessageType       = "message_types.PanicMessage"
	ActorName              = "_actor_"
	GlobalChanSize         = 10
	DefaultActorPollSize   = 5
	SentimentActorPool     = "sentimentActorPool"
	AggregationActorPool   = "aggregationActorPool"

	ClusterDatabaseAddress         = "mongodb+srv://snowtoslow:123qweASD@magiccluster.ccit5.mongodb.net/tweet-db?retryWrites=true&w=majority"
	DatabaseName                   = "tweet-db"
	TweetsCollection               = "tweets"
	UserCollection                 = "users"
	TickerInterval                 = 200 * time.Millisecond
	UserModel                      = "models.User"
	EngagementRatio                = "*models.EngagementRation"
	SentimentAnalysis              = "*models.SentimentAnalysis"
	NameOfFileToGenerateRandomUuid = "/dev/urandom"
)
