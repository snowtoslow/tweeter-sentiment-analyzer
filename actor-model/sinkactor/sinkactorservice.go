package sinkactor

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"tweeter-sentiment-analyzer/actor-model/actorabstraction"
	"tweeter-sentiment-analyzer/actor-model/actorregistry"
	"tweeter-sentiment-analyzer/constants"
)

func NewSinkActor(actorName string) actorabstraction.IActor {
	chanToRecv := make(chan interface{}, constants.GlobalChanSize)
	mapForStorage := make(map[string][]interface{})
	sinkActor := &SinkActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToRecv,
		},
		sinkBuffer: mapForStorage,
	}

	go sinkActor.ActorLoop()

	(*actorregistry.MyActorRegistry)["sinkActor"] = sinkActor

	return sinkActor
}

func (sinkActor *SinkActor) ActorLoop() {
	defer close(sinkActor.ActorProps.ChanToReceiveData)
	ticker := time.NewTicker(constants.TickerInterval)
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(constants.ClusterDatabaseAddress))
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		select {
		case action := <-sinkActor.ActorProps.ChanToReceiveData:
			if fmt.Sprintf("%T", action) == "models.User" {
				sinkActor.sinkBuffer["users"] = append(sinkActor.sinkBuffer["users"], action)
			} else if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct || fmt.Sprintf("%T", action) == constants.RetweetedStatus {
				sinkActor.sinkBuffer["tweets"] = append(sinkActor.sinkBuffer["tweets"], action)
			}

			if len(sinkActor.sinkBuffer["users"])+len(sinkActor.sinkBuffer["tweets"]) == 128 {
				log.Println("full buffer!")
				if err = sinkActor.insertAndClear(mongoClient); err != nil {
					log.Fatal(err)
				}
				ticker.Reset(constants.TickerInterval)
			}
		case <-ticker.C:
			log.Println("after 200ms:", len(sinkActor.sinkBuffer["users"])+len(sinkActor.sinkBuffer["tweets"]))
			if err = sinkActor.insertAndClear(mongoClient); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (sinkActor *SinkActor) insertAndClear(mongoClient *mongo.Client) (errorOccurredInInsert error) {
	myArray := sinkActor.sinkBuffer
	for k, v := range myArray {
		go func(k string, v []interface{}) {
			if k == "users" {
				if _, err := mongoClient.Database(constants.DatabaseName).Collection(constants.UserCollection).InsertMany(context.Background(), v); err != nil {
					errorOccurredInInsert = err
					log.Fatal("FATAL ERROR INSERTING USERS:", err)
				}
			} else if k == "tweets" {
				if _, err := mongoClient.Database(constants.DatabaseName).Collection(constants.TweetsCollection).InsertMany(context.Background(), v); err != nil {
					errorOccurredInInsert = err
					log.Fatal("FATAL ERROR INSERTING TWEETS:", err)
				}
			}
		}(k, v)
	}

	sinkActor.sinkBuffer["users"] = sinkActor.sinkBuffer["users"][:0]
	sinkActor.sinkBuffer["tweets"] = sinkActor.sinkBuffer["tweets"][:0]
	//sinkActor.sinkBuffer = make(map[string][]interface{})
	return
}

func (sinkActor *SinkActor) SendMessage(msg interface{}) {
	sinkActor.ActorProps.ChanToReceiveData <- msg
}
