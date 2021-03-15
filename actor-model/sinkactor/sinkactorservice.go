package sinkactor

import (
	"context"
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
	sinkActor := &SinkActor{
		ActorProps: actorabstraction.AbstractActor{
			Identity:          actorName + constants.ActorName,
			ChanToReceiveData: chanToRecv,
		},
	}

	go sinkActor.ActorLoop()

	(*actorregistry.MyActorRegistry)["sinkActor"] = sinkActor

	return sinkActor
}

func (sinkActor *SinkActor) ActorLoop() {
	defer close(sinkActor.ActorProps.ChanToReceiveData)
	ticker := time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case action := <-sinkActor.ActorProps.ChanToReceiveData:
			sinkActor.SinkBuffer = append(sinkActor.SinkBuffer, action)
			if len(sinkActor.SinkBuffer) == 128 {
				log.Println("FULL BUFFER:", len(sinkActor.SinkBuffer))
				myArray := sinkActor.SinkBuffer
				go func() {
					if err := sinkActor.insertDataInDb(myArray); err != nil {
						log.Fatal("Error inserting 200 ms!", err)
					}
				}()
				sinkActor.SinkBuffer = sinkActor.SinkBuffer[:0]
				ticker.Reset(200 * time.Millisecond)
			}
		case <-ticker.C:
			log.Println("after 200ms:", len(sinkActor.SinkBuffer))
			myArray := sinkActor.SinkBuffer
			go func() {
				if err := sinkActor.insertDataInDb(myArray); err != nil {
					log.Fatal("Error inserting 200 ms!", err)
				}
			}()
			sinkActor.SinkBuffer = sinkActor.SinkBuffer[:0]
		}
	}
}

func (sinkActor *SinkActor) insertAndClear(array []interface{}) error {
	var myeer error
	go func() {
		if err := sinkActor.insertDataInDb(array); err != nil {
			myeer = err
			log.Fatal("Error inserting 200 ms!", err)
		}
	}()
	sinkActor.SinkBuffer = sinkActor.SinkBuffer[:0]
	return myeer
}

func (sinkActor *SinkActor) insertDataInDb(array []interface{}) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(constants.ClusterDatabaseAddress))
	if err != nil {
		return
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			return
		}
	}()

	if _, err = mongoClient.Database(constants.DatabaseName).Collection(constants.TweetsCollection).InsertMany(ctx, array); err != nil {
		return
	}

	return nil
}

func (sinkActor *SinkActor) SendMessage(msg interface{}) {
	sinkActor.ActorProps.ChanToReceiveData <- msg
}
