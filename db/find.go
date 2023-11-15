package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FrosTiK-SD/mongik/constants"
	"github.com/allegro/bigcache/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getKey(collectionName string, operation string, query bson.M) string {
	return fmt.Sprintf("%s | %s | %v", collectionName, constants.DB_FINDONE, query)
}

// Should be called after every write operation on a cluster
func resetCache(cacheClient bigcache.BigCache, clusterName string) {
	DBCacheReset(&cacheClient, clusterName)
}

func FindOne[Result any](mongoClient *mongo.Client, cacheClient *bigcache.BigCache,db string, collectionName string, query bson.M, result *Result, noCache bool) {
	key := getKey(collectionName, constants.DB_FINDONE, query)
	var resultBytes []byte
	var resultInterface map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes, _ := cacheClient.Get(key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retreiving DB call from the cache with cache key ", key)
			return
		}
	}

	// Query to DB
	mongoClient.Database(db).Collection(collectionName).FindOne(context.Background(), query).Decode(&resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	// Set to cache
	resultBytes, _ = json.Marshal(result)
	if err := DBCacheSet(cacheClient, key, resultBytes); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}
}

func Find[Result any]( mongoClient *mongo.Client, cacheClient *bigcache.BigCache, db string,collectionName string, query bson.M, noCache bool) ([]Result, error) {
	key := getKey(collectionName, constants.DB_FIND, query)
	var resultBytes []byte
	var result []Result
	var resultInterface []map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes, _ := cacheClient.Get(key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retreiving DB call from the cache with cache key ", key)
			return result, nil
		}
	}

	// Query to DB
	fmt.Println("Queriying the DB")
	cursor, err := mongoClient.Database(db).Collection(collectionName).Find(context.Background(), query)
	if err != nil {
		return nil, err
	}
	cursor.All(context.Background(), &resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	// Set to cache
	resultBytes, _ = json.Marshal(result)
	if err := DBCacheSet(cacheClient, key, resultBytes); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}

	return result, nil
}
