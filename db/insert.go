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
func InsertMany[Result any](mongikClient *mongik.Mongik, db string, collectionName string, doc []interface{}) (*mongo.InsertManyResult, error) {	
	// Query to DB
	docIds, err := mongikClient.MongoClient.Database(db).Collection(collectionName).InsertMany(context.Background(), doc)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docIds, err
}