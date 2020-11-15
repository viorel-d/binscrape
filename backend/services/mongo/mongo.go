package mongo

import (
	"context"
	"log"
	"time"

	"github.com/viorel-d/binscrape/backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func newClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.MongoDBURI).SetAuth(
		options.Credential{
			AuthSource: config.MongoAuthSource,
			Username:   config.MongoUsername,
			Password:   config.MongoPassword,
		})
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func ping(client *mongo.Client) (err error) {
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	log.Println("Received pong...")

	return
}

func GetClient() *mongo.Client {
	checkClient := func(client *mongo.Client) {
		err := ping(client)
		if err != nil {
			log.Printf("err: %#v\n", err)
			log.Fatalln("Couldn't connect to Mongo")
		}
	}
	if mongoClient != nil {
		checkClient(mongoClient)
		return mongoClient
	}
	mongoClient = newClient()
	checkClient(mongoClient)
	log.Println("Successfully connected to Mongo")

	return mongoClient
}

func GetCollection(name string) *mongo.Collection {
	client := GetClient()
	db := client.Database(config.MongoDBName)
	col := db.Collection(name)

	return col
}

func CreateIndex(
	collection *mongo.Collection,
	keysMap *bson.M,
	options *options.IndexOptions,
) (bool, error) {
	index := mongo.IndexModel{
		Keys:    *keysMap,
		Options: options,
	}
	if _, err := collection.Indexes().CreateOne(context.TODO(), index); err != nil {
		return false, err
	}

	return true, nil
}

// CreateUniqueIndex attempts to create an unique index using the given keysMap
func CreateUniqueIndex(collection *mongo.Collection, keysMap *bson.M) (bool, error) {
	return CreateIndex(collection, keysMap, options.Index().SetUnique(true))
}

// ListIndexes returns the indexes for the given collection
func ListIndexes(collection *mongo.Collection) []bson.M {
	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := collection.Indexes().List(context.TODO(), opts)
	if err != nil {
		log.Fatalln(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatalln(err)
	}

	return results
}

func EnsureIndex(collection *mongo.Collection, keyName string) bool {
	colIndexes := ListIndexes(collection)

	for _, colIndex := range colIndexes {
		colIndexKey, err := AsBsonM(colIndex["key"])
		if err != nil {
			log.Fatalln(err)
		}
		index := colIndexKey[keyName]
		if index != nil {
			return false
		}
	}

	return true
}

func InsertMany(collection *mongo.Collection, items *[]interface{}) {
	var itemsToInsert []interface{}
	for _, item := range *items {
		b, err := bson.Marshal(item)
		if err != nil {
			log.Fatalln(err)
		}
		var bsonMapItem bson.M
		err = bson.Unmarshal(b, &bsonMapItem)
		if err != nil {
			log.Fatalln(err)
		}
		itemsToInsert = append(itemsToInsert, bsonMapItem)
	}
	result, err := collection.InsertMany(context.TODO(), itemsToInsert)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully inserted: %#v\n", result.InsertedIDs)
}

func Find(collection *mongo.Collection, opts *options.FindOptions) []interface{} {
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		log.Fatalln(err)
	}
	defer cursor.Close(context.TODO())
	var items []interface{}
	for cursor.Next(context.TODO()) {
		var item interface{}
		err := cursor.Decode(&item)
		if err != nil {
			log.Fatalln(err)
		}

		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalln(err)
	}

	return items
}

func FindOne(
	collection *mongo.Collection,
	id string,
	opts *options.FindOneOptions,
) (result interface{}, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}, opts).Decode(&result)

	return
}
