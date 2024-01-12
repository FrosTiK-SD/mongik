package mongik

import (
	"context"

	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOne[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	// Query to DB
	result, err := mongikClient.MongoClient.Database(db).Collection(collectionName).UpdateOne(context.Background(), filter, update, opts...)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return result, err
}

func UpdateMany[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	// Query to DB
	result, err := mongikClient.MongoClient.Database(db).Collection(collectionName).UpdateMany(context.Background(), filter, update, opts...)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return result, err
}
