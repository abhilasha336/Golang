package driver

import (
	"context"
	"fmt"

	"logger/internal/entities"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectDB initializes db
func ConnectDB(cfg entities.Database) (*mongo.Database, error) {
	ctx := context.Background()
	_string := prepareConnectionString(cfg)
	fmt.Println("[Mongo] connecting in.....", _string)

	clientOptions := options.Client().ApplyURI(_string)

	// Connect to the MongoDB server
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.DATABASE)
	return db, nil

}

func prepareConnectionString(cfg entities.Database) string {
	var str string
	if cfg.Driver != "mongodb" {
		str = fmt.Sprintf("%s://%s:%s@%s", cfg.Driver, cfg.User, cfg.Password, cfg.Host)
	} else {
		str = fmt.Sprintf("%s://%s:%s@%s:%d", cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port)
	}
	return str
}
