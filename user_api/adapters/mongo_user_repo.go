package adapters

import (
	"context"
	"fmt"
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "users"
	databaseName   = "users"
)

type MongoRepository struct {
	*mongoConnection
	db *mongo.Database
}

func NewMongoRepository() (*MongoRepository, error) {
	m := &MongoRepository{
		mongoConnection: &mongoConnection{},
	}
	loadConfigError := m.mongoConnection.loadConfig()
	if loadConfigError != nil {
		return m, loadConfigError
	}
	m.initialized = true
	err := m.verifyConnection(context.Background())
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m MongoRepository) Save(ctx context.Context, user *domain.User) error {
	// Get the client
	client, err := m.Connection(ctx)
	if err != nil {
		return err
	}

	// Run command
	_, err = client.Database(databaseName).Collection(collectionName).InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoRepository) UpdateByID(ctx context.Context, id string, user *domain.User) error {
	// Get the client
	client, err := m.Connection(ctx)
	if err != nil {
		return err
	}
	// Run command
	upResult, err := client.Database(databaseName).Collection(collectionName).UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", user}})
	if err != nil {
		return err
	}

	if upResult.MatchedCount == 0 {
		return fmt.Errorf("no user with Id %s found", id)
	}

	return nil
}

func (m MongoRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	// Get the client
	client, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}
	var resp domain.User
	// Run command
	err = client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.D{{"_id", id}}).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (m MongoRepository) DeleteById(ctx context.Context, id string) error {
	// Get the client
	client, err := m.Connection(ctx)
	if err != nil {
		return err
	}

	// Run command
	_, err = client.Database(databaseName).Collection(collectionName).DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return nil

}
