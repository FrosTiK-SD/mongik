package constants

import mongik "github.com/FrosTiK-SD/mongik/models"

const (
	BIGCACHE = "BIGCACHE"
	REDIS    = "REDIS"
)

var DEFAULT_REDIS_CONFIG = &mongik.RedisConfig{
	URI: "localhost:6379",
	DBPassword: "",
	DBIndex: 0,
}