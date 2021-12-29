package adapters

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

type mongoConnection struct {
	ConnectionURL string
	Username      string
	Password      string
	client        *mongo.Client
	clientOptions *options.ClientOptions
	initialized   bool
	sync.Mutex
}

func (c *mongoConnection) loadConfig() error {
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
	opts, err := c.makeClientOpts()
	if err != nil {
		return err
	}

	c.clientOptions = opts
	return nil
}
func (c *mongoConnection) makeClientOpts() (*options.ClientOptions, error) {
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
func (c *mongoConnection) createClient(ctx context.Context) (client *mongo.Client, err error) {
	if !c.initialized {
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
func (c *mongoConnection) Close() error {
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

func (c *mongoConnection) verifyConnection(ctx context.Context) error {
	client, err := c.createClient(ctx)
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

func (c *mongoConnection) Connection(ctx context.Context) (*mongo.Client, error) {
	if !c.initialized {
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

	client, err := c.createClient(ctx)
	if err != nil {
		return nil, err
	}
	c.client = client
	return c.client, nil
}
