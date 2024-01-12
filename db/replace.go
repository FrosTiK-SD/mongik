package mongik

import (
	"context"

	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ReplaceOne[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, filter bson.M, doc Doc, noCache bool, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	// Query to DB
	docId, err := mongikClient.MongoClient.Database(db).Collection(collectionName).ReplaceOne(context.Background(), filter, doc, opts...)

	DBCacheReset(mongikClient.CacheClient, collectionName)
	return docId, err
}