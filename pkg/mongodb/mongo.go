package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/ahd99/urlshortner/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const URLShortner_DB_Name = "urlshortner"
const Req_Log_Collection_Name = "reqlog"

var logger1 logger.Logger
var mongoClient *mongo.Client

func InitMongo(logg logger.Logger) {
	logger1 = logg
	clientOpts := options.Client()
	clientOpts.ApplyURI("mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000")
	mongoClient1, err := mongo.NewClient(clientOpts)
	if err != nil {
		logger1.Fatal("Error in mongo.NewClient", logger.String("err", err.Error()))
	}

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	err = mongoClient1.Connect(ctx)
	if err != nil {
		logger1.Fatal("Error connecting to mongo",  logger.String("err", err.Error()))
	}

	err = mongoClient1.Ping(ctx, readpref.Primary())
	if err != nil {
		logger1.Fatal("mongo ping error", logger.String("err", err.Error()))
	}

	databases, err := mongoClient1.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		logger1.Fatal("Error getting database names.", logger.String("err", err.Error()))
	}
	logger1.Info(fmt.Sprintf("databases: %v", databases))

	mongoClient = mongoClient1

}

func InsertReqLog(key string, url string, srcip string, t time.Time) {
	reqLogCollection := mongoClient.Database(URLShortner_DB_Name).Collection(Req_Log_Collection_Name)
	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)

	s, err := reqLogCollection.InsertOne(ctx, bson.D{
		{Key: "ukey", Value: key}, 
		{"url", url},
		{"ip", srcip},
		{"reqTime", t},
	})
	if err != nil {
		logger1.Error("Error inserting req log", 
		logger.String("key", key), 
		logger.String("url", url), 
		logger.String("ip", srcip), 
		logger.String("time", t.String()),
		logger.String("err", err.Error()))
	} else {
		logger1.Debug("Ereq log instered", 
		logger.String("key", key), 
		logger.String("url", url), 
		logger.String("ip", srcip), 
		logger.String("time", t.String()),
		logger.String("_id", fmt.Sprintf("%v", s.InsertedID)),
		)
	}
}

func CloseMongo() {
	ctx, _ := context.WithTimeout(context.Background(), 2 * time.Second)
	mongoClient.Disconnect(ctx)
}