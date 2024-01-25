package mongik

import (
	"context"

	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertOne[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, doc Doc, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	// Query to DB
	docId, err := mongikClient.MongoClient.Database(db).Collection(collectionName).InsertOne(context.Background(), doc, opts...)

	DBCacheReset(mongikClient, collectionName)
	return docId, err
}

func InsertMany[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, docs []Doc, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	// Convert []struct to []interface
	docsInterface := make([]interface{}, len(docs))
	for index := range docs {
		docsInterface[index] = docs[index]
	}
	// Query to DB
	docIds, err := mongikClient.MongoClient.Database(db).Collection(collectionName).InsertMany(context.Background(), docsInterface, opts...)

	DBCacheReset(mongikClient, collectionName)
	return docIds, err
}
