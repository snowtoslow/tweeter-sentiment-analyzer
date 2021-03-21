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
	"tweeter-sentiment-analyzer/models"
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
			if fmt.Sprintf("%T", action) == constants.UserModel {
				log.Println(action.(models.User).UniqueId)
				sinkActor.sinkBuffer[constants.UserCollection] = append(sinkActor.sinkBuffer[constants.UserCollection], action)
			} else if fmt.Sprintf("%T", action) == constants.JsonNameOfStruct || fmt.Sprintf("%T", action) == constants.RetweetedStatus {
				sinkActor.sinkBuffer[constants.TweetsCollection] = append(sinkActor.sinkBuffer[constants.TweetsCollection], action)
			}

			if len(sinkActor.sinkBuffer[constants.UserCollection])+len(sinkActor.sinkBuffer[constants.TweetsCollection]) == 128 {
				log.Println("full buffer!")
				if err = sinkActor.insertAndClear(mongoClient); err != nil {
					log.Fatal(err)
				}
				ticker.Reset(constants.TickerInterval)
			}
		case <-ticker.C:
			//log.Println("after 200ms:", len(sinkActor.sinkBuffer[constants.UserCollection])+len(sinkActor.sinkBuffer[constants.TweetsCollection]))
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
			if k == constants.UserCollection {
				if _, err := mongoClient.Database(constants.DatabaseName).Collection(constants.UserCollection).InsertMany(context.Background(), v); err != nil {
					errorOccurredInInsert = err
					log.Fatal("FATAL ERROR INSERTING USERS:", err)
				}
			} else if k == constants.TweetsCollection {
				if _, err := mongoClient.Database(constants.DatabaseName).Collection(constants.TweetsCollection).InsertMany(context.Background(), v); err != nil {
					errorOccurredInInsert = err
					log.Fatal("FATAL ERROR INSERTING TWEETS:", err)
				}
			}
		}(k, v)
	}

	/*sinkActor.sinkBuffer[constants.UserCollection] = sinkActor.sinkBuffer[constants.UserCollection][:0]
	sinkActor.sinkBuffer[constants.TweetsCollection] = sinkActor.sinkBuffer[constants.TweetsCollection][:0]*/
	sinkActor.sinkBuffer = make(map[string][]interface{})
	return
}

func (sinkActor *SinkActor) SendMessage(msg interface{}) {
	sinkActor.ActorProps.ChanToReceiveData <- msg
}
