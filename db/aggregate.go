package mongik

import (
	"context"
	"fmt"

	"github.com/FrosTiK-SD/mongik/constants"
	mongik "github.com/FrosTiK-SD/mongik/models"
)

func Aggregate[Result any](mongikClient *mongik.Mongik, db string, collectionName string, pipeline interface{}, noCache bool) ([]Result, error) {
	var option interface{}
	key := getKey(collectionName, constants.DB_AGGREGATE, pipeline, option)
	var resultBytes []byte
	var result []Result
	var resultInterface []map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes, _ := mongikClient.CacheClient.Get(key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retreiving DB call from the cache with cache key ", key)
			return result, nil
		}
	}

	// Query to DB
	fmt.Println("Queriying the DB")
	cursor, err := mongikClient.MongoClient.Database(db).Collection(collectionName).Aggregate(context.Background(), pipeline)
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

func AggregateOne[Result any](mongikClient *mongik.Mongik, db string, collectionName string, pipeline interface{}, noCache bool) (Result, error) {
	var option interface{}
	key := getKey(collectionName, constants.DB_AGGREGATEONE, pipeline, option)
	var resultBytes []byte
	var result Result
	var resultInterface []map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes, _ := mongikClient.CacheClient.Get(key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retreiving DB call from the cache with cache key ", key)
			return result, nil
		}
	}

	// Query to DB
	fmt.Println("Queriying the DB")
	cursor, err := mongikClient.MongoClient.Database(db).Collection(collectionName).Aggregate(context.Background(), pipeline)
	if err != nil {
		return result, err
	}
	cursor.All(context.Background(), &resultInterface)

	if len(resultInterface) == 0 {
		return result, constants.ERROR_NO_DOCS
	}

	resultBody, _ := json.Marshal(resultInterface[0])
	json.Unmarshal(resultBody, &result)

	// Set to cache
	resultBytes, _ = json.Marshal(result)
	if err := DBCacheSet(mongikClient.CacheClient, key, resultBytes); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}

	return result, nil
}
