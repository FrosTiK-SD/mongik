package mongik

import (
	"context"
	"log"
	"time"

	mongik "github.com/FrosTiK-SD/mongik/models"
	"github.com/allegro/bigcache/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func NewClient(mongoURL string, cachingDuration time.Duration) *mongik.Mongik {
	// Initialize BigCache
	cacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(cachingDuration))
	
	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB: %v\n", err)
	} else {
		log.Println("Connected to MongoDB")
	}

	return &mongik.Mongik{
		MongoClient: mongoClient,
		CacheClient: cacheClient,
	}
}