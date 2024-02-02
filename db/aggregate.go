package mongik

import (
	"context"
	"fmt"

	"github.com/FrosTiK-SD/mongik/constants"
	mongik "github.com/FrosTiK-SD/mongik/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Aggregate[Result any](mongikClient *mongik.Mongik, db string, collectionName string, pipeline []bson.M, noCache bool, opts ...*options.AggregateOptions) ([]Result, error) {
	key := getKey(collectionName, constants.DB_AGGREGATE, pipeline, opts)
	var result []Result
	var resultInterface []map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes := DBCacheFetch(mongikClient, key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			return result, nil
		}
	}

	CacheLog(mongikClient, fmt.Sprintf("Querying the DB\n KEY: %s", key))

	// Query to DB
	cursor, err := mongikClient.MongoClient.Database(db).Collection(collectionName).Aggregate(context.Background(), pipeline, opts...)
	if err != nil {
		return nil, err
	}
	cursor.All(context.Background(), &resultInterface)

	resultBody, _ := json.Marshal(resultInterface)
	json.Unmarshal(resultBody, &result)

	// Parsing lookup collection from pipeline
	lookupCollections := getLookupCollections(pipeline)

	// Set to cache
	DBCacheSet(mongikClient, key, result, lookupCollections...)

	return result, nil
}

func AggregateOne[Result any](mongikClient *mongik.Mongik, db string, collectionName string, pipeline []bson.M, noCache bool, opts ...*options.AggregateOptions) (Result, error) {
	key := getKey(collectionName, constants.DB_AGGREGATEONE, pipeline, opts)
	var result Result
	var resultInterface []map[string]interface{}

	// First Check if it is present in the cache
	if !noCache {
		resultBytes := DBCacheFetch(mongikClient, key)
		if err := json.Unmarshal(resultBytes, &result); err == nil {
			return result, nil
		}
	}

	CacheLog(mongikClient, fmt.Sprintf("Querying the DB\n KEY: %s", key))

	// Query to DB
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

	// Parsing lookup collection from pipeline
	lookupCollections := getLookupCollections(pipeline)

	// Set to cache
	DBCacheSet(mongikClient, key, result, lookupCollections...)

	return result, nil
}
