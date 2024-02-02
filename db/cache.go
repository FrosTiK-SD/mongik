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
	var keyStoreBytes []byte

	// Get the list of keys
	if mongikClient.Config.Client == constants.BIGCACHE {
		keyStoreBytes, _ = mongikClient.CacheClient.Get(constants.KEY_STORE)
	} else if mongikClient.Config.Client == constants.REDIS {
		keyStoreBytes, _ = mongikClient.RedisClient.Get(context.Background(), constants.KEY_STORE).Bytes()
	}

	if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
		keyStore = make(map[string][]string)
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
	keyStoreBytes, _ = json.Marshal(keyStore)
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

	CacheLog(mongikClient, fmt.Sprintf("Keystore set: %s", keyStore))

	return nil
}

func DBCacheReset(mongikClient *mongik.Mongik, clusterName string) {
	var keyStore map[string][]string
	var keyStoreBytes []byte

	// Get the list of keys
	if mongikClient.Config.Client == constants.BIGCACHE {
		keyStoreBytes, _ = mongikClient.CacheClient.Get(constants.KEY_STORE)
	} else if mongikClient.Config.Client == constants.REDIS {
		keyStoreBytes, _ = mongikClient.RedisClient.Get(context.Background(), constants.KEY_STORE).Bytes()
	}

	if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
		keyStore = make(map[string][]string)
	}

	// Delete all the keys in the cluster
	if mongikClient.Config.Client == constants.BIGCACHE {
		for _, key := range keyStore[clusterName] {
			mongikClient.CacheClient.Delete(key)
		}
	} else if mongikClient.Config.Client == constants.REDIS {
		for _, key := range keyStore[clusterName] {
			mongikClient.RedisClient.Del(context.Background(), key)
		}
	}

	keyStore[clusterName] = []string{}

	// Set the key store
	keyStoreBytes, _ = json.Marshal(keyStore)

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
		CacheLog(mongikClient, fmt.Sprintf("Retrieved data from the cache of the key: %s", key))
		return resultBytes
	} else if mongikClient.Config.Client == constants.REDIS {
		resultBytes, _ := mongikClient.RedisClient.Get(context.Background(), key).Bytes()
		CacheLog(mongikClient, fmt.Sprintf("Retrieved data from the cache of the key: %s", key))
		return resultBytes
	}
	return nil
}
