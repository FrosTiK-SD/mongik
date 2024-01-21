package mongik

import (
	"context"
	"fmt"
	"strings"

	"github.com/FrosTiK-SD/mongik/constants"
	mongik "github.com/FrosTiK-SD/mongik/models"
)

func getDBClusterFromKey(key string) string {
	// We store DB cache in the format CLUSTER_NAME | OPERATION | QUERY
	return strings.Split(key, " | ")[0]
}

func DBCacheSet(mongikClient *mongik.Mongik, key string, value interface{}, lookupCollections ...string) error {
	ctx := context.Background()
	var keyStore map[string][]string

	// Get the list of keys
	if mongikClient.Config.Client == constants.BIGCACHE {
		keyStoreBytes, _ := mongikClient.CacheClient.Get(constants.KEY_STORE)
		if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
			keyStore = make(map[string][]string)
		}
	} else if mongikClient.Config.Client == constants.REDIS {
		keyStoreResult, _ := mongikClient.RedisClient.HGetAll(ctx, constants.KEY_STORE).Result()
		keyStoreBytes, err := json.Marshal(keyStoreResult)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
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

	fmt.Println("--------- Keystore set")
	fmt.Println(keyStore)
	fmt.Println("--------- Keystore set end")

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
		err := mongikClient.RedisClient.Set(ctx, constants.KEY_STORE, keyStoreBytes, mongikClient.Config.TTL).Err()
		if err != nil {
			return err
		}
		return mongikClient.RedisClient.Set(ctx, key, valueBytes, mongikClient.Config.TTL).Err()
	}

	return nil
}

func DBCacheReset(mongikClient *mongik.Mongik, clusterName string) {
	ctx := context.Background()
	var keyStore map[string][]string

	// Get the list of keys
	if mongikClient.Config.Client == constants.BIGCACHE {
		keyStoreBytes, _ := mongikClient.CacheClient.Get(constants.KEY_STORE)
		if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
			keyStore = make(map[string][]string)
		}
	} else if mongikClient.Config.Client == constants.REDIS {
		keyStoreResult := mongikClient.RedisClient.HGetAll(ctx, constants.KEY_STORE)
		keyStoreBytes, _ := json.Marshal(keyStoreResult)
		fmt.Println(keyStore)
		if err := json.Unmarshal(keyStoreBytes, &keyStore); err != nil {
			keyStore = make(map[string][]string)
		}
	}

	// Delete all the keys in the cluster
	for _, key := range keyStore[clusterName] {
		mongikClient.RedisClient.Del(ctx, key)
	}
	keyStore[clusterName] = []string{}

	// Set the key store
	if mongikClient.Config.Client == constants.BIGCACHE {
		keyStoreBytes, _ := json.Marshal(keyStore)
		mongikClient.CacheClient.Set(constants.KEY_STORE, keyStoreBytes)
	} else if mongikClient.Config.Client == constants.REDIS {
		mongikClient.RedisClient.HSet(ctx, constants.KEY_STORE, keyStore)
	}
}

func DBCacheFetch(mongikClient *mongik.Mongik, key string) []byte {
	ctx := context.Background()
	if mongikClient.Config.Client == constants.BIGCACHE {
		resultBytes, _ := mongikClient.CacheClient.Get(key)
		return resultBytes
	} else if mongikClient.Config.Client == constants.REDIS {
		exists, _ := mongikClient.RedisClient.Exists(ctx, key).Result()
		if exists == 0 {
			return nil
		}
		result, _ := mongikClient.RedisClient.HGetAll(ctx, key).Result()
		if resultBytes, err := json.Marshal(result); err == nil {
			return resultBytes
		}
	}
	return nil
}
