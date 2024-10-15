package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Config   Config
	Client   *mongo.Client
	Database *mongo.Database
}

func (d *Database) Connect() {
	opts := options.Client()
	uri := d.Config.BuildUri()

	opts.ApplyURI(uri)

	fmt.Println("[info] Connecting to MongoDB")
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Printf("[error] Failed to connect to MongoDB Database: %s", uri)
		panic(1) // panic here as this is a fatal error
	}

	d.Database = client.Database("mtgjson")
	d.Client = client
}

func (d Database) Disconnect() {
	d.Health() // this will throw an fatal error when

	fmt.Println("[info] Disconnecting from MongoDB") // these will be replaced with proper logging in a future PR
	err := d.Client.Disconnect(context.Background())
	if err != nil {
		fmt.Println("[error] Failed to disconnect from MongoDB: ", err)
		panic(1)
	}
}

func (d Database) Health() {
	err := d.Client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("[error] Failed to ping MongoDB")
		panic(1)
	}
}

func (d Database) Find(collection string, query bson.D, model interface{}) any {
	coll := d.Database.Collection(collection)

	var results interface{}
	err := coll.FindOne(context.TODO(), query).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}

	bytes, err := bson.Marshal(results)
	if err != nil {
		fmt.Println("[error] Failed to marshal results:", err)
	}

	err2 := bson.Unmarshal(bytes, model)
	if err2 != nil {
		fmt.Println("[error] Failed to unmarshal results:", err2)
	}

	return model
}
