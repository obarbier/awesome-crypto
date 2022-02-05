package common

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"
)

type MongoConnection struct {
	ConnectionURL string
	Username      string
	Password      string
	client        *mongo.Client
	clientOptions *options.ClientOptions
	Initialized   bool
	sync.Mutex
}

func (c *MongoConnection) LoadConfig() error {
	// FIXME
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	err := viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	opts, err := c.MakeClientOpts()
	if err != nil {
		return err
	}
	// TODO: Add index similar to https://github.com/go-session/mongo/blob/master/mongo.go#L35
	c.clientOptions = opts
	return nil
}
func (c *MongoConnection) MakeClientOpts() (*options.ClientOptions, error) {
	//writeOpts, err := c.getWriteConcern()
	//if err != nil {
	//	return nil, err
	//}
	//
	//authOpts, err := c.getTLSAuth()
	//if err != nil {
	//	return nil, err
	//}
	//
	//timeoutOpts, err := c.timeoutOpts()
	//if err != nil {
	//	return nil, err
	//}

	opts := options.MergeClientOptions()
	return opts, nil
}
func (c *MongoConnection) CreateClient(ctx context.Context) (client *mongo.Client, err error) {
	if !c.Initialized {
		return nil, fmt.Errorf("failed to create client: connection producer is not initialized")
	}
	if c.clientOptions == nil {
		return nil, fmt.Errorf("missing client options")
	}
	client, err = mongo.Connect(ctx, options.MergeClientOptions(options.Client().ApplyURI(c.ConnectionURL), c.clientOptions))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Close terminates the database connection.
func (c *MongoConnection) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := c.client.Disconnect(ctx); err != nil {
			return err
		}
	}

	c.client = nil

	return nil
}

func (c *MongoConnection) VerifyConnection(ctx context.Context) error {
	client, err := c.CreateClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify connection: %w", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		_ = client.Disconnect(ctx) // Try to prevent any sort of resource leak
		return fmt.Errorf("failed to verify connection: %w", err)
	}
	c.client = client
	return nil
}

func (c *MongoConnection) Connection(ctx context.Context) (*mongo.Client, error) {
	if !c.Initialized {
		return nil, fmt.Errorf("database client is not initialized")
	}

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.client != nil {
		if err := c.client.Ping(ctx, readpref.Primary()); err == nil {
			return c.client, nil
		}
		// Ignore error on purpose since we want to re-create a session
		_ = c.client.Disconnect(ctx)
	}

	client, err := c.CreateClient(ctx)
	if err != nil {
		return nil, err
	}
	c.client = client
	return c.client, nil
}
