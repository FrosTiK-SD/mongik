package mongik

import (
	"context"

	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteOne(mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	// Query to DB
	docId, err := mongikClient.MongoClient.Database(db).Collection(collectionName).DeleteOne(context.Background(), query, opts...)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docId, err
}

func DeleteMany(mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	// Query to DB
	docIds, err := mongikClient.MongoClient.Database(db).Collection(collectionName).DeleteMany(context.Background(), query, opts...)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docIds, err
}