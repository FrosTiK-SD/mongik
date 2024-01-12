package mongik

import (
	"context"

	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOne[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, filter bson.M, doc Doc, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	// Query to DB
	docId, err := mongikClient.MongoClient.Database(db).Collection(collectionName).UpdateOne(context.Background(), filter, doc, opts...)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docId, err
}

func UpdateMany[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, filter bson.M, docs []Doc, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	// Convert []struct to []interface
	docsInterface := make([]interface{}, len(docs))
	for index := range docs {
		docsInterface[index] = docs[index]
	}
	// Query to DB
	docIds, err := mongikClient.MongoClient.Database(db).Collection(collectionName).UpdateMany(context.Background(), filter, docsInterface, opts...)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docIds, err
}

