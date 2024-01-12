package mongik

import (
	"context"

	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOne[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, doc Doc) (*mongo.InsertOneResult, error) {
	// Query to DB
	docId, err := mongikClient.MongoClient.Database(db).Collection(collectionName).InsertOne(context.Background(), doc)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docId, err
}

func InsertMany[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, docs []Doc) (*mongo.InsertManyResult, error) {
	// Convert []struct to []interface
	docsInterface := make([]interface{}, len(docs))
	for index := range docs {
		docsInterface[index] = docs[index]
	}
	// Query to DB
	docIds, err := mongikClient.MongoClient.Database(db).Collection(collectionName).InsertMany(context.Background(), docsInterface)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docIds, err
}
