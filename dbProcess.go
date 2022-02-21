package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbInsert(date string, title string, descriptor string, primary_subject string, act_link string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.0.105:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("testing").Collection("numbers")
	res, err := collection.InsertOne(ctx, bson.D{{"date", date}, {"title", title}, {"descriptor", descriptor}, {"primary_subject", primary_subject}, {"act_link", act_link}})
	if err != nil {
		log.Errorln(err)
		return
	}
	id := res.InsertedID
	log.Println("mongo id ", id)
}
