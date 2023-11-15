package mongik

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/allegro/bigcache/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongik struct {
	MongoClient *mongo.Client
	CacheClient *bigcache.BigCache
}

func NewClient(mongoURL string, cachingDuration time.Duration) *Mongik {
	// Initialize BigCache
	cacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(cachingDuration))
	
	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv(mongoURL)).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB: %v\n", err)
	} else {
		log.Println("Connected to MongoDB")
	}

	return &Mongik{
		MongoClient: mongoClient,
		CacheClient: cacheClient,
	}
}