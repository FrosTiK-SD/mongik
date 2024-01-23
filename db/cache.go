package mongik

import (
	"context"
	"fmt"
	"strings"

	"github.com/FrosTiK-SD/mongik/constants"
	mongik "github.com/FrosTiK-SD/mongik/models"
)

func getDBClusterFromKey(key string) string {

	// We store DB cache in the format CLUSTER_NAME | OPERATION | QUERY | OPTIONS
	return strings.Split(key, " | ")[0]
}

func DBCacheSet(mongikClient *mongik.Mongik, key string, value interface{}, lookupCollections ...string) error {
	var keyStore map[string][]string

	// Get the list of keys
	if mongikClient.Config.Client == constants.BIGCACHE {

		keyStoreBytes, _ := mongikClient.CacheClient.Get(constants.KEY_STORE)
		if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
			keyStore = make(map[string][]string)
		}

	} else if mongikClient.Config.Client == constants.REDIS {

		keyStoreResult, _ := mongikClient.RedisClient.Get(context.Background(), constants.KEY_STORE).Result()
		if err := json.UnmarshalFromString(keyStoreResult, &keyStore); err != nil {
			keyStore = make(map[string][]string)
		}

	}

	clusterName := getDBClusterFromKey(key)

	// Add it to the cluster set
	keyStore[clusterName] = append(keyStore[clusterName], key)
	if lookupCollections != nil {
		for _, collection := range lookupCollections {
			keyStore[collection] = append(keyStore[collection], key)
		}
	}

	// Set the key store
	keyStoreBytes, _ := json.Marshal(keyStore)
	valueBytes, _ := json.Marshal(value)

	if mongikClient.Config.Client == constants.BIGCACHE {
		err := mongikClient.CacheClient.Set(constants.KEY_STORE, keyStoreBytes)
		if err != nil {
			return err
		}
		return mongikClient.CacheClient.Set(key, valueBytes)
	} else if mongikClient.Config.Client == constants.REDIS {
		err := mongikClient.RedisClient.Set(context.Background(), constants.KEY_STORE, keyStoreBytes, mongikClient.Config.TTL).Err()
		if err != nil {
			return err
		}
		return mongikClient.RedisClient.Set(context.Background(), key, valueBytes, mongikClient.Config.TTL).Err()
	}

	fmt.Println("Keystore set: ", keyStore)

	return nil
}

func DBCacheReset(mongikClient *mongik.Mongik, clusterName string) {
	var keyStore map[string][]string

	// Get the list of keys
	if mongikClient.Config.Client == constants.BIGCACHE {

		keyStoreBytes, _ := mongikClient.CacheClient.Get(constants.KEY_STORE)
		if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
			keyStore = make(map[string][]string)
		}

	} else if mongikClient.Config.Client == constants.REDIS {

		keyStoreResult, _ := mongikClient.RedisClient.Get(context.Background(), constants.KEY_STORE).Result()
		if err := json.UnmarshalFromString(keyStoreResult, &keyStore); err != nil {
			keyStore = make(map[string][]string)
		}

	}

	// Delete all the keys in the cluster
	for _, key := range keyStore[clusterName] {
		mongikClient.RedisClient.Del(context.Background(), key)
	}
	keyStore[clusterName] = []string{}

	// Set the key store
	keyStoreBytes, _ := json.Marshal(keyStore)

	if mongikClient.Config.Client == constants.BIGCACHE {

		err := mongikClient.CacheClient.Set(constants.KEY_STORE, keyStoreBytes)
		if err != nil {
			fmt.Println("Error in setting Keystore: ", err)
		}

	} else if mongikClient.Config.Client == constants.REDIS {

		err := mongikClient.RedisClient.Set(context.Background(), constants.KEY_STORE, keyStoreBytes, mongikClient.Config.TTL).Err()
		if err != nil {
			fmt.Println("Error in setting Keystore: ", err)
		}

	}
}

func DBCacheFetch(mongikClient *mongik.Mongik, key string) []byte {
	
	// Fetch from Cache
	if mongikClient.Config.Client == constants.BIGCACHE {

		resultBytes, _ := mongikClient.CacheClient.Get(key)
		return resultBytes

	} else if mongikClient.Config.Client == constants.REDIS {

		result, _ := mongikClient.RedisClient.Get(context.Background(), key).Result()

		// This is done as result is string type and marshalling directly to JSON throws error
		var resultInterface map[string]interface{}
		if err := json.UnmarshalFromString(result, &resultInterface); err != nil {
			return nil
		}
		if resultBytes, err := json.Marshal(resultInterface); err == nil {
			return resultBytes
		}

	}
	return nil
}
