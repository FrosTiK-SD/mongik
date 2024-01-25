# Mongik

A simple MongoDB wrapper for Golang MongoDB driver. It does not try to replace the native MongoDB driver but simply act as a wrapper with in-built caching and optimised query parsing.

Provides two caching options, either using `BigCache` or using `Redis`.

## Why use Mongik?

If you want to reduce DB calls to MongoDB to save DB costs and you are getting lots of `READ` requests and less `WRITE` requests then Mongik will be the right choice for you. This is a wrapper over the MongoDB driver so you can use the native driver anytime.

`Redis` version of Mongik is also available for use, if you have a horizontally scaled server.

## Why to not use Mongik?

If you have lots of `WRITE` requests and less `READ` requests, Mongik may not be a life-saver for you. You can still use Mongik but you will not have any significant performance gain over the native MongoDB driver as it will use the native MongoDB driver itself for the `WRITE` operations as there is no point in caching in this scenario.

If you have a scaled server (Horizontally scaled specifically) then use `Redis` version of `Mongik`. [BETA]

## How to get started?

It is pretty simple!

### Installation

```.go
go get github.com/FrosTiK-SD/mongik
```

### Initialisation

`BigCache` version

```.go
package main

import (
    "time"

    mongik "github.com/FrosTiK-SD/mongik"
    models "github.com/FrosTiK-SD/mongik/models"
)

func main() {
    mongikConfig := &models.config{
        Client: "BIGCACHE",
        TTL: time.Hour,
    }
    mongikClient := mongik.NewClient(os.Getenv(constants.DB), mongikConfig)
}
```

`Redis` version

```.go
package main

import (
    "time"

    mongik "github.com/FrosTiK-SD/mongik"
    models "github.com/FrosTiK-SD/mongik/models"
)

func main() {
    mongikConfig := &models.config{
        Client: "REDIS",
        RedisConfig: &models.RedisConfig{
            URI: "localhost:6379",          
            DBPassword: "",
            DBIndex: 0,
        }                                   // Default config if RedisConfig left empty
        TTL: time.Hour,
        FallbackToDefault: true,            // If true, will default to BigCache version if Redis throws error
    }
    mongikClient := mongik.NewClient(os.Getenv(constants.DB), mongikConfig)
}
```

No error handling required! Any error while connecting to Mongo will be managed by the `MongikClient`

### Parameters of NewClient function
| Parameter No | Name | Type | Usage |
| ------------ | ---- | ---- | ----- |
| 1 | MONGO_CONNECTION_STRING | `string` | The `MongoDB` connection string `mongodb+srv://.....` |
| 2 | MONGIK_CONFIG | `Config` from `models/mongik.go` | Specify client version and other configurations |

Below is the `Config` struct.

```.go
type Config struct {
	Client string                       // Specify "BIGCACHE" or "REDIS"
	RedisConfig *RedisConfig
	TTL time.Duration
	FallbackToDefault bool
}

type RedisConfig struct {
	URI string                          // Redis server addr
	DBPassword string
	DBIndex int
}
``` 

The function returns a `MongikClient`.

```.go
type Mongik struct {
    MongoClient *mongo.Client
	CacheClient *bigcache.BigCache      // Populated in BigCache version, else empty
	RedisClient *redis.Client           // Populated in Redis version, else empty
	Config *Config
}
```

You can also use the individual clients for more granular control.

### All done!

Now you can check out the `db` folder to see the methods exported. Replace the DB calls in your code with the exported functions of `Mongik` to enjoy caching and enhanced parsing!
