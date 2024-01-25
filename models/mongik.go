package mongik

import (
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongik struct {
	MongoClient *mongo.Client
	CacheClient *bigcache.BigCache
	RedisClient *redis.Client
	Config      *Config
}

type Config struct {
	Client            string
	RedisConfig       *RedisConfig
	TTL               time.Duration
	FallbackToDefault bool
}

type RedisConfig struct {
	URI      string
	Username string
	Password string
	DBIndex  int
}
