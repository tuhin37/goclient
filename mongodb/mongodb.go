package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ----------------- client operations ------------------

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

// select collection
func (mongodb *MongodbClient) Collection(collName string) mongoCollection {
	return mongodb.Collections[collName]
}

// ----------------- custom datatypes ------------------
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoInsertedOne struct {
	result *mongo.InsertOneResult
}

type mongoInsertedMany struct {
	result *mongo.InsertManyResult
}

type mongoSingleResult struct {
	result *mongo.SingleResult
}

type mongoCursor struct {
	cursor *mongo.Cursor
}

type mongoUpdateResult struct {
	result *mongo.UpdateResult
}

type mongoDeleteResult struct {
	result *mongo.DeleteResult
}

type MongodbClient struct {
	URI         string
	Databse     string
	Collections map[string]mongoCollection
	Client      *mongo.Client
}

// ----------------- collection operations ------------------
func (coll mongoCollection) Aggregate(pipeline primitive.A) mongoCursor {
	cursor, _ := coll.coll.Aggregate(context.Background(), pipeline)
	return mongoCursor{cursor: cursor}
}

func (coll mongoCollection) InsertOne(document any) (*mongoInsertedOne, error) {
	result, err := coll.coll.InsertOne(context.Background(), document)
	if err != nil {
		return &mongoInsertedOne{}, err
	}
	return &mongoInsertedOne{result: result}, nil
}

func (coll mongoCollection) InsertMany(documents []any) (*mongoInsertedMany, error) {
	result, err := coll.coll.InsertMany(context.Background(), documents)
	if err != nil {
		return &mongoInsertedMany{}, err
	}
	return &mongoInsertedMany{result: result}, nil
}

func (coll mongoCollection) FindOne(filter any) *mongoSingleResult {
	result := coll.coll.FindOne(context.Background(), filter)
	return &mongoSingleResult{result: result}
}

func (coll mongoCollection) Find(filter any) (*mongoCursor, error) {
	cursor, err := coll.coll.Find(context.Background(), filter)
	if err != nil {
		return &mongoCursor{}, err
	}
	return &mongoCursor{cursor: cursor}, nil
}

func (coll mongoCollection) UpdateOne(filter any, update any) (*mongoUpdateResult, error) {
	result, err := coll.coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return &mongoUpdateResult{}, err
	}
	return &mongoUpdateResult{result: result}, nil
}

func (coll mongoCollection) UpdateMany(filter any, update any) (*mongoUpdateResult, error) {
	result, err := coll.coll.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return &mongoUpdateResult{}, err
	}
	return &mongoUpdateResult{result: result}, nil
}

func (coll mongoCollection) DeleteOne(filter any) (*mongoDeleteResult, error) {
	result, err := coll.coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return &mongoDeleteResult{}, err
	}
	return &mongoDeleteResult{result: result}, nil
}

func (coll mongoCollection) DeleteMany(filter any) (*mongoDeleteResult, error) {
	result, err := coll.coll.DeleteMany(context.Background(), filter)
	if err != nil {
		return &mongoDeleteResult{}, err
	}
	return &mongoDeleteResult{result: result}, nil
}

// ----------------- decode ------------------
func (i *mongoInsertedOne) GetID() string {
	return i.result.InsertedID.(primitive.ObjectID).Hex()
}

func (i *mongoInsertedMany) GetIDs() []string {
	var ids []string
	for _, id := range i.result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID).Hex())
	}
	return ids
}

func (r *mongoSingleResult) Decode(v any) error {
	return r.result.Decode(v)
}

func (r *mongoCursor) GetDocuments() ([]interface{}, error) {
	defer r.cursor.Close(context.Background())

	var results []interface{}

	for r.cursor.Next(context.Background()) {
		var v interface{}
		err := r.cursor.Decode(&v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document: %s", err)
		}

		results = append(results, v)
	}

	if err := r.cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %s", err)
	}
	return results, nil
}

func (m *mongoUpdateResult) MatchedCount() int64 {
	return m.result.MatchedCount
}

func (m *mongoUpdateResult) ModifiedCount() int64 {
	return m.result.ModifiedCount
}

func (m *mongoUpdateResult) UpsertedCount() int64 {
	return m.result.UpsertedCount
}

func (m *mongoUpdateResult) UpsertedID() string {
	return m.result.UpsertedID.(primitive.ObjectID).Hex()
}

func (m *mongoDeleteResult) DeletedCount() int64 {
	return m.result.DeletedCount
}

/*
------------------------------- example --------------------------------
package main

import (
	"fmt"

	"github.com/tuhin37/goclient/mongodb"
)

var mongo *mongodb.MongodbClient
var err error

func init() {
	mongo, err = mongodb.NewMongodbClient("mongodb+srv://mongoroot:Hn4Wp1LJnUgFBmiC@cluster0.rnjm0.mongodb.net", "drag", []string{"batches", "tasks"})
	if err != nil {
		panic(err)
	}
}

func main() {
	defer mongo.Close()
	mongo.Collection("batches").InsertOne(map[string]interface{}{"foo": "bar"})
	fmt.Println("done")
}


*/
