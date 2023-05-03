package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCursor struct {
	cursor *mongo.Cursor
}

func (cursor mongoCursor) Close() {
	cursor.cursor.Close(context.Background())
}

// METHOD: decode the cursor into a struct
func (cursor mongoCursor) Decode(v any) error {
	var results []interface{}

	// load cursor into map
	if err := cursor.cursor.All(context.Background(), &results); err != nil {
		return err
	}

	// map to bson []byte
	bytes, err := bson.Marshal(results)
	if err != nil {
		return err
	}

	// []byte to struct
	if err := bson.Unmarshal(bytes, v); err != nil {
		return err
	}
	return nil
}

type mongoCollection struct {
	coll *mongo.Collection
}

type MongodbClient struct {
	URI         string
	Databse     string
	Collections map[string]mongoCollection
	Client      *mongo.Client
}

// CONSTRUCTOR: creates a new mongodb client
func NewMongodbClient(uri string, database string, collections []string) (*MongodbClient, error) {
	mongodb := MongodbClient{}
	var err error

	opts := options.Client().ApplyURI(uri)                    // client options
	mongodb.Client, err = mongo.Connect(context.TODO(), opts) // connect to mongodb
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// ATP: connection successful
	// create & load all the collections-handlers
	mongodb.Databse = database
	mongodb.Collections = make(map[string]mongoCollection)
	for _, collName := range collections {
		mongodb.Collections[collName] = mongoCollection{coll: (*mongo.Collection)(mongodb.Client.Database(database).Collection(collName))}
	}

	return &mongodb, nil
}

// METHOD: close the mongodb connection. Defer this function after creating a new client
func (mongodb *MongodbClient) Close() {
	mongodb.Client.Disconnect(context.Background())
}

// METHOD: returns a collection handler, given the collection name
func (mongodb *MongodbClient) Collection(collName string) mongoCollection {
	return mongodb.Collections[collName]
}

// METHOD: on a collection and returns a cursor
func (coll mongoCollection) Aggregate(pipeline primitive.A) mongoCursor {
	cursor, _ := coll.coll.Aggregate(context.Background(), pipeline)
	return mongoCursor{cursor: cursor}
}
