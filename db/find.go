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
	key := getKey(collectionName, constants.DB_FINDONE, query, opts)
	var resultInterface map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes := DBCacheFetch(mongikClient, key)
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
	if err := DBCacheSet(mongikClient, key, result); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}
}

func Find[Result any](mongikClient *mongik.Mongik, db string, collectionName string, query bson.M, noCache bool, opts ...*options.FindOptions) ([]Result, error) {
	key := getKey(collectionName, constants.DB_FIND, query, opts)
	var result []Result
	var resultInterface []map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes := DBCacheFetch(mongikClient, key)
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
	if err := DBCacheSet(mongikClient, key, result); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}

	return result, nil
}
