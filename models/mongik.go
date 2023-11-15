package models

import (
	"github.com/allegro/bigcache/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongik struct {
	MongoClient *mongo.Client
	CacheClient *bigcache.BigCache
}
