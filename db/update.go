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

	DBCacheReset(mongikClient, collectionName)
	return result, err
}

func UpdateMany[Doc any](mongikClient *mongik.Mongik, db string, collectionName string, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	// Query to DB
	result, err := mongikClient.MongoClient.Database(db).Collection(collectionName).UpdateMany(context.Background(), filter, update, opts...)

	DBCacheReset(mongikClient, collectionName)
	return result, err
}

func FindOneAndUpdate[Result any](mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, update bson.M, opts ...*options.FindOneAndUpdateOptions) Result {
	var result Result
	var resultInterface map[string]interface{}

	mongikClient.MongoClient.Database(db).Collection(collectionName).FindOneAndUpdate(context.Background(), query, update, opts...).Decode(&resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	DBCacheReset(mongikClient, collectionName)

	return result
}
