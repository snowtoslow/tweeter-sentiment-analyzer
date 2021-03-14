package sinkactor

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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
				sinkActor.SinkBuffer = sinkActor.SinkBuffer[:0]
			}
		case <-ticker.C:
			log.Println("after 200ms:", len(sinkActor.SinkBuffer))
			sinkActor.SinkBuffer = sinkActor.SinkBuffer[:0]
		}
	}
}

// WithTransactionExample is an example of using the Session.WithTransaction function.
func (sinkActor *SinkActor) WithTransactionExample() {
	ctx := context.Background()

	clientOpts := options.Client().ApplyURI(constants.ClusterDatabaseAddress)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Prereq: Create collections.
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(1*time.Second))
	wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)
	fooColl := client.Database(constants.DatabaseName).Collection(constants.TweetsCollection, wcMajorityCollectionOpts)

	// Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
		// transaction.

		res, err := fooColl.InsertMany(sessCtx, sinkActor.SinkBuffer)
		if err != nil {
			log.Println("ERR:", err)
			return nil, err
		}

		return res, nil
	}

	// Step 2: Start a session and run the callback using WithTransaction.
	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %v\n", result)
}

func (sinkActor *SinkActor) InsertDataInDb() (err error) {
	/*if err := sinkActor.InsertDataInDb();err!=nil{
		log.Fatal("Error inserting!",err)
	}*/
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(constants.ClusterDatabaseAddress))
	if err != nil {
		return err
	}

	if _, err = mongoClient.Database(constants.DatabaseName).Collection(constants.TweetsCollection).InsertMany(context.Background(), sinkActor.SinkBuffer); err != nil {
		return err
	}

	return nil
}

func (sinkActor *SinkActor) SendMessage(msg interface{}) {
	sinkActor.ActorProps.ChanToReceiveData <- msg
}
