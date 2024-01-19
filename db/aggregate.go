package mongik

import (
	"context"
	"fmt"
	"strings"

	"github.com/FrosTiK-SD/mongik/constants"
	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Aggregate[Result any](mongikClient *mongik.Mongik, db string, collectionName string, pipeline interface{}, noCache bool, opts ...*options.AggregateOptions) ([]Result, error) {
	key := getKey(collectionName, constants.DB_AGGREGATE, pipeline, opts)
	var resultBytes []byte
	var result []Result
	var resultInterface []map[string]interface{}

	// Parsing lookup collection from pipeline
	var lookupCollection string = " "
	pipe := fmt.Sprintf("%v", pipeline)
	pipeSplit := strings.Split(pipe, "$lookup")
	if len(pipeSplit) > 1 {
		pipeSplit2 := strings.Split(pipeSplit[1], " ")
		for _, tag := range pipeSplit2 {
			if strings.Contains(tag, "from:") {
				res := strings.Split(tag, "from:")
				lookupCollection = res[1]
				break
			}
		}
	}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes := DBCacheFetch(mongikClient.CacheClient, key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retrieving DB call from the cache with cache key ", key)
			return result, nil
		}
	}

	// Query to DB
	fmt.Println("Querying the DB")
	cursor, err := mongikClient.MongoClient.Database(db).Collection(collectionName).Aggregate(context.Background(), pipeline, opts...)
	if err != nil {
		return nil, err
	}
	cursor.All(context.Background(), &resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	// Set to cache
	resultBytes, _ = json.Marshal(result)
	if err := DBCacheSet(mongikClient.CacheClient, key, resultBytes, lookupCollection); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}

	return result, nil
}

func AggregateOne[Result any](mongikClient *mongik.Mongik, db string, collectionName string, pipeline interface{}, noCache bool, opts ...*options.AggregateOptions) (Result, error) {
	key := getKey(collectionName, constants.DB_AGGREGATEONE, pipeline, opts)
	var resultBytes []byte
	var result Result
	var resultInterface []map[string]interface{}

	// Parsing lookup collection from pipeline
	var lookupCollection string = " "
	pipe := fmt.Sprintf("%v", pipeline)
	pipeSplit := strings.Split(pipe, "$lookup")
	if len(pipeSplit) > 1 {
		pipeSplit2 := strings.Split(pipeSplit[1], " ")
		for _, tag := range pipeSplit2 {
			if strings.Contains(tag, "from:") {
				res := strings.Split(tag, "from:")
				lookupCollection = res[1]
				break
			}
		}
	}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes := DBCacheFetch(mongikClient.CacheClient, key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			fmt.Println("Retrieving DB call from the cache with cache key ", key)
			return result, nil
		}
	}

	// Query to DB
	fmt.Println("Querying the DB")
	cursor, err := mongikClient.MongoClient.Database(db).Collection(collectionName).Aggregate(context.Background(), pipeline, opts...)
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
	if err := DBCacheSet(mongikClient.CacheClient, key, resultBytes, lookupCollection); err == nil {
		fmt.Println("Successfully set DB call in cache with key ", key)
	}

	return result, nil
}
