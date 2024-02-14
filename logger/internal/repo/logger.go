package repo

import (
	"context"
	"log"
	"logger/internal/consts"
	"logger/internal/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
LoggerRepo represents a repository responsible for interacting with a database
to retrieve log data. It encapsulates database-specific operations related to logs.
*/
type LoggerRepo struct {
	db *mongo.Database
}

// LoggerRepoImply specifies the contract for interacting with the LoggerRepo repository.
type LoggerRepoImply interface {
	AddLog(ctx context.Context, log entities.Log) error
	GetLogs(ctx context.Context, filters bson.M, page, limit int32) ([]*entities.Log, int64, error)
}

/*
NewLoggerRepo creates a new instance of the LoggerRepo repository, initializing it
with the provided MongoDB database connection.
*/
func NewLoggerRepo(db *mongo.Database) LoggerRepoImply {
	return &LoggerRepo{db: db}
}

// AddLog adds a new log entry in the database.
func (logger *LoggerRepo) AddLog(ctx context.Context, log entities.Log) error {
	collection := logger.db.Collection(consts.CollectionLogs)
	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	return nil
}

/*
GetLogs retrieves log data from the database based on provided filters, pagination,
and returns the logs as well as the total number of matching records.
*/
func (logger *LoggerRepo) GetLogs(ctx context.Context, filters bson.M, page, limit int32) ([]*entities.Log, int64, error) {

	collection := logger.db.Collection(consts.CollectionLogs)
	l, skip := int64(limit), int64(page*limit-limit)
	cursor, err := collection.Find(ctx, filters, options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(l).SetSkip(skip))
	if err != nil {
		return nil, 0, err
	}
	var logs []*entities.Log
	for cursor.Next(ctx) {
		var currentLog entities.Log
		if err := cursor.Decode(&currentLog); err != nil {
			log.Println(err)
		}
		

		logs = append(logs, &currentLog)
	}
	totalRecords, err := collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}
	return logs, totalRecords, nil
}
