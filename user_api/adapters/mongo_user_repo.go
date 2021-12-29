package adapters

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"strings"
)

type MongoRepository struct {
	*mongoConnection
	db *mongo.Database
}

const (
	collectionName = "users"
	databaseName   = "users"
)

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

func (m MongoRepository) Save(ctx context.Context, firstName, lastName, userId, passwordHash string) (domain.User, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return domain.User{}, err
	}
	cmd := createUserCommand{
		Collection: collectionName,
		User: []domain.User{{
			Id:           id.String(),
			FirstName:    firstName,
			LastName:     lastName,
			UserId:       userId,
			PasswordHash: passwordHash,
		}},
	}
	if err := m.runCommandWithRetry(ctx, databaseName, cmd); err != nil {
		return domain.User{}, err
	}

	resp := domain.User{}
	return resp, nil
}

func (m MongoRepository) Update(ctx context.Context, id, firstName, lastName, userId, passwordHash string) (domain.User, error) {
	cmd := updateUserCommand{
		Collection: collectionName,
		Updates: []updateStatement{
			{
				Query:   bson.D{{"_id", id}},
				Updates: bson.D{{"$set", bson.D{{"firstName", firstName}, {"lastName", lastName}, {"userId", userId}, {"passwordHash", passwordHash}}}},
			},
		},
	}
	if err := m.runCommandWithRetry(ctx, databaseName, cmd); err != nil {
		return domain.User{}, err
	}

	resp := domain.User{}
	return resp, nil
}

func (m MongoRepository) Get(ctx context.Context, userId string) (domain.User, error) {
	cmd := findUserCommand{
		Collection: collectionName,
		Filter:     bson.D{{"_id", userId}},
	}
	if err := m.runCommandWithRetry(ctx, databaseName, cmd); err != nil {
		return domain.User{}, err
	}

	resp := domain.User{}
	return resp, nil
}

func (m MongoRepository) Delete(ctx context.Context, userId string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoRepository) runCommandWithRetry(ctx context.Context, db string, cmd interface{}) error {
	// Get the client
	client, err := m.Connection(ctx)
	if err != nil {
		return err
	}

	// Run command
	cursor, err := client.Database(db).RunCommandCursor(ctx, cmd, nil) // Error check on the first attempt
	var results []bson.D
	// check for errors in the conversion
	if cursor != nil {
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		// display the documents retrieved
		fmt.Println("displaying all results from the search query")
		for _, result := range results {
			log.Printf("%+v\n", result)
		}

	}

	switch {
	case err == nil:
		return nil
	case err == io.EOF, strings.Contains(err.Error(), "EOF"):
		// Call getConnection to reset and retry query if we get an EOF error on first attempt.
		client, err = m.Connection(ctx)
		if err != nil {
			return err
		}
		cursor, err = client.Database(db).RunCommandCursor(ctx, cmd, nil)
		if err != nil {
			return err
		}
	case strings.EqualFold(err.Error(), "cursor should be an embedded document but is of BSON type invalid"):
		log.Printf("error due to RunCommandCursor %v", err)
	default:
		return err
	}

	return nil
}
