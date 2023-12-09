# Mongik

A simple MongoDB warpper for Golang MongoDB driver. It no way tries to replace the native MongoDB driver and just acts like a wrapper with in-built caching and optimised query parsing.

## Why use Mongik ?

Suppose you want to reduce DB calls to MongoDB to save costing you are getting lots of `READ` requests and less `WRITE`, Mongik will be the right choice for you. This is just a wrapper over MongoDB driver so you can use the native driver anytime according to your preference.

## Why not use Mongik ?

If you have a use-case where there are lots of `WRITE` requests and less `READ` requests, Mongik may not be a life-saver for you. You can still use Mongik but you will not have any significant performance gain over the native MongoDB driver as it will use the native MongoDB driver itself for the operations can there is no point in caching in this scenario.

If you have a scaled server (Horizontally scaled specifically) then Mongik will not be a good choice for you **if you are using the `BigCache` version of Mongik.** In such a usecase you have to use the `Redis` version of `Mongik` [BETA]

## How to get started ?

Its actually very simple

Lets first install it

```.go
go get github.com/FrosTiK-SD/mongik
```

Initialize it

```.go
package main

import (
    mongik "github.com/FrosTiK-SD/mongik"
)

func main() {
    mongikClient := mongik.NewClient(os.Getenv(constants.DB), time.Hour)
}
```

Its that simple. All the error in connecting to Mongo are managed by the `MongikClient` itself

| Parameter No | Name | Type | Usage |
| ------------ | ---- | ---- | ----- |
| 1 | MONGO_CONNECTION_STRING | `string` | The `MongoDB` connection string `mongodb+srv://.....` |
| 2 | CACHE_DURATION | `time.Duration` | The duration for which the DB call will be cached |

It returns a `MongikClient`.

```.go
type Mongik struct {
    MongoClient *mongo.Client
    CacheClient *bigcache.BigCache
}
```

You can use the individual clients also for more granular control but we will not talk about here as documentation of that can be found in their respective docs.

Now you can check out the `db` folder to check out the method exported and just replace the DB calls in your code with the exported functions of `Mongik` to enjoy caching and enhanced parsing.
