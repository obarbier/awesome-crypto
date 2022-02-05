package adapters

import (
	"context"
	common "github.com/obarbier/awesome-crypto/internal"
	"github.com/obarbier/awesome-crypto/pkg/crypto/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "keys"
	databaseName   = "keys"
)

type MongoKeyRepository struct {
	*common.MongoConnection
	db *mongo.Database
}

var _ domain.KeyRepository = MongoKeyRepository{}

func NewMongoKeyRepository() (*MongoKeyRepository, error) {
	m := &MongoKeyRepository{
		MongoConnection: &common.MongoConnection{},
	}
	loadConfigError := m.MongoConnection.LoadConfig()
	if loadConfigError != nil {
		return m, loadConfigError
	}
	m.Initialized = true
	err := m.VerifyConnection(context.Background())
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m MongoKeyRepository) Save(ctx context.Context, k domain.KeyEntity) error {
	// Get the client
	client, err := m.Connection(ctx)
	if err != nil {
		return err
	}

	// Run command
	_, err = client.Database(databaseName).Collection(collectionName).InsertOne(ctx, k)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoKeyRepository) Get(ctx context.Context, id string) (k *domain.KeyEntity, err error) {
	// Get the client
	client, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}
	var resp domain.KeyEntity
	// Run command
	err = client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.D{{"_id", id}}).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
