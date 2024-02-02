package mongik

import (
	"context"
	"log"

	"github.com/FrosTiK-SD/mongik/constants"
	mongik "github.com/FrosTiK-SD/mongik/models"
	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(mongoURL string, config *mongik.Config) *mongik.Mongik {

	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB: %v\n", err)
	} else {
		log.Println("Connected to MongoDB")
	}

	// Initialize BigCache anyway
	cacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(config.TTL))

	// Check for caching mode
	if config.Client == constants.REDIS {
		// Check for default redisConfig
		if config.RedisConfig == nil {
			config.RedisConfig = constants.DEFAULT_REDIS_CONFIG
		}
		// Initialize Redis
		redisClient := redis.NewClient(&redis.Options{
			Addr:     config.RedisConfig.URI,
			Username: config.RedisConfig.Username,
			Password: config.RedisConfig.Password,
			DB:       config.RedisConfig.DBIndex,
		})
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			if config.FallbackToDefault == true {
				return &mongik.Mongik{
					MongoClient: mongoClient,
					CacheClient: cacheClient,
					Config: &mongik.Config{
						Debug:             config.Debug,
						Client:            constants.BIGCACHE,
						TTL:               config.TTL,
						FallbackToDefault: config.FallbackToDefault,
					},
				}
			} else {
				log.Fatalf("Unable to Connect to Redis: %v\n", err)
			}
		}
		return &mongik.Mongik{
			MongoClient: mongoClient,
			RedisClient: redisClient,
			CacheClient: cacheClient,
			Config:      config,
		}
	}

	return &mongik.Mongik{
		MongoClient: mongoClient,
		CacheClient: cacheClient,
		Config:      config,
	}
}
