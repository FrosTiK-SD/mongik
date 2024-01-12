package mongik

import (
	"context"
	"fmt"

	"github.com/FrosTiK-SD/mongik/constants"
	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne[Result any](mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, result *Result, noCache bool, opts ...*options.FindOneOptions) {
	var option interface{}
	key := getKey(collectionName, constants.DB_FINDONE, query, option)
	var resultBytes []byte
	var resultInterface map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes, _ := mongikClient.CacheClient.Get(key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retrieving DB call from the cache with cache key ", key)
			return
		}
	}

	// Query to DB
	mongikClient.MongoClient.Database(db).Collection(collectionName).FindOne(context.Background(), query, opts...).Decode(&resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	// Set to cache
	resultBytes, _ = json.Marshal(result)
	if err := DBCacheSet(mongikClient.CacheClient, key, resultBytes); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}
}

func Find[Result any](mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, noCache bool, opts ...*options.FindOptions) ([]Result, error) {
	key := getKey(collectionName, constants.DB_FIND, query, opts)
	var resultBytes []byte
	var result []Result
	var resultInterface []map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes, _ := mongikClient.CacheClient.Get(key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retrieving DB call from the cache with cache key ", key)
			return result, nil
		}
	}

	// Query to DB
	fmt.Println("Querying the DB")
	cursor, err := mongikClient.MongoClient.Database(db).Collection(collectionName).Find(context.Background(), query, opts...)
	if err != nil {
		return nil, err
	}
	cursor.All(context.Background(), &resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	// Set to cache
	resultBytes, _ = json.Marshal(result)
	if err := DBCacheSet(mongikClient.CacheClient, key, resultBytes); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}

	return result, nil
}

func FindOneAndUpdate[Result any](mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, update bson.M, noCache bool, opts ...*options.FindOneAndUpdateOptions) Result {
	var result Result
	var resultInterface map[string]interface{}

	fmt.Println("Querying the DB")
	mongikClient.MongoClient.Database(db).Collection(collectionName).FindOneAndUpdate(context.Background(), query, update, opts...).Decode(&resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	DBCacheReset(mongikClient.CacheClient, collectionName)

	return result
}

func FindOneAndReplace[Result any](mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, replace bson.M, noCache bool, opts ...*options.FindOneAndReplaceOptions) Result {
	var result Result
	var resultInterface map[string]interface{}

	fmt.Println("Querying the DB")
	mongikClient.MongoClient.Database(db).Collection(collectionName).FindOneAndReplace(context.Background(), query, replace, opts...).Decode(&resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	DBCacheReset(mongikClient.CacheClient, collectionName)

	return result
}
