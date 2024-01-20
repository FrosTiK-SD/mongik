package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FrosTiK-SD/mongik/constants"
	db "github.com/FrosTiK-SD/mongik/db"
	mongik "github.com/FrosTiK-SD/mongik/models"
	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Test struct {
	Id primitive.ObjectID `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
	Val int64 `bson:"val" json:"val"`
}

func main() {
	mongikClient := NewClient("mongodb://localhost:27017/", &mongik.Config{
		Client: "REDIS",
		TTL: time.Hour,
	})
	var res Test
	db.FindOne[Test](mongikClient, "test", "test", bson.M{
		"_id": "65a2cc7c86f33746a787bcb9",
	}, &res, false)
	fmt.Println("---res", res)
	db.FindOne[Test](mongikClient, "test", "test", bson.M{
		"_id": "65a2cc7c86f33746a787bcb9",
	}, &res, false)
	fmt.Println("---res", res)
	db.FindOne[Test](mongikClient, "test", "test", bson.M{
		"_id": "65a2cc7c86f33746a787bcb9",
	}, &res, false)
	fmt.Println("---res", res)
}

func NewClient(mongoURL string, config *mongik.Config) *mongik.Mongik {
	ctx := context.Background()

	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB: %v\n", err)
	} else {
		log.Println("Connected to MongoDB")
	}

	// Check for caching mode
	if config.Client == constants.REDIS {
		// Check for default redisConfig
		if config.RedisConfig == nil {
			config.RedisConfig = constants.DEFAULT_REDIS_CONFIG
		}
		// Initialize Redis
		redisClient := redis.NewClient(&redis.Options{
			Addr:     config.RedisConfig.URI,
			Password: config.RedisConfig.DBPassword,
			DB:       config.RedisConfig.DBIndex,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			if config.FallbackToDefault == true {
				// Initialize BigCache
				cacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(config.TTL))
				return &mongik.Mongik{
					MongoClient: mongoClient,
					CacheClient: cacheClient,
					Config:      config,
				}
			} else {
				log.Fatalf("Unable to Connect to Redis: %v\n", err)
			}
		}
		return &mongik.Mongik{
			MongoClient: mongoClient,
			RedisClient: redisClient,
			Config:      config,
		}

	} else if config.Client == constants.BIGCACHE {
		// Initialize BigCache
		cacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(config.TTL))
		return &mongik.Mongik{
			MongoClient: mongoClient,
			CacheClient: cacheClient,
			Config:      config,
		}
	}

	return &mongik.Mongik{
		MongoClient: mongoClient,
		Config:      config,
	}
}
