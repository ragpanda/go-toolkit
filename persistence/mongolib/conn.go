package mongolib

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ragpanda/go-toolkit/bizerr"
	"github.com/ragpanda/go-toolkit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Name     string `json:"Name" yaml:"Name"`
	MongoURI string `json:"MongoURI" yaml:"MongoURI"`
}

type MongoDBPoolConfig struct {
	Config []*MongoConfig `json:"Config" yaml:"Config"`
}

type MongoDBPool struct {
	config        MongoDBPoolConfig
	connectionMap sync.Map
	mu            sync.Mutex
}

var globalPool *MongoDBPool

func SetGlobalPool(ctx context.Context, pool *MongoDBPool) {
	globalPool = pool
}
func GetGlobalPool(ctx context.Context) *MongoDBPool {
	return globalPool
}

func NewGlobalPool(ctx context.Context, config MongoDBPoolConfig) (*MongoDBPool, error) {
	var err error
	globalPool, err = NewMongoDBPool(ctx, config)
	return globalPool, err
}

func NewMongoDBPool(ctx context.Context, config MongoDBPoolConfig) (*MongoDBPool, error) {
	self := &MongoDBPool{config: config}
	err := self.initConnections(ctx)
	return self, err
}

func (pool *MongoDBPool) initConnections(ctx context.Context) error {
	for _, cfg := range pool.config.Config {
		client, err := pool.createConnection(ctx, cfg.MongoURI)
		if err != nil {
			// Log the error or handle it as appropriate for your application
			log.Error(ctx, "Failed to create connection", err)
			return err
		}
		pool.connectionMap.Store(cfg.Name, client)
	}
	return nil
}

func (pool *MongoDBPool) createConnection(ctx context.Context, uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (pool *MongoDBPool) GetConnection(ctx context.Context, name string) (*mongo.Client, error) {
	if client, ok := pool.connectionMap.Load(name); ok {
		return client.(*mongo.Client), nil
	}
	return nil, bizerr.ErrNotFound.WithMessage("db connection not found").WithStack(ctx)
}

func (pool *MongoDBPool) SetConfig(ctx context.Context, cfg *MongoConfig) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Create a new connection
	newClient, err := pool.createConnection(ctx, cfg.MongoURI)
	if err != nil {
		return err
	}

	// Check if the connection already exists
	if oldClient, ok := pool.connectionMap.Load(cfg.Name); ok {
		// Close the old connection
		if err := oldClient.(*mongo.Client).Disconnect(context.Background()); err != nil {
			return err
		}
	}

	// Store the new connection
	pool.connectionMap.Store(cfg.Name, newClient)

	// Update the config
	for i, c := range pool.config.Config {
		if c.Name == cfg.Name {
			pool.config.Config[i] = cfg
			return nil
		}
	}
	pool.config.Config = append(pool.config.Config, cfg)

	return nil
}

func (pool *MongoDBPool) CloseAll(ctx context.Context) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.connectionMap.Range(func(key, value interface{}) bool {
		if err := value.(*mongo.Client).Disconnect(ctx); err != nil {
			return false
		}
		pool.connectionMap.Delete(key)
		return true
	})

	return nil
}

func (pool *MongoDBPool) CloseConnection(ctx context.Context, name string) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if client, ok := pool.connectionMap.Load(name); ok {
		if err := client.(*mongo.Client).Disconnect(ctx); err != nil {
			return err
		}
		pool.connectionMap.Delete(name)
		return nil
	}

	return bizerr.ErrNotFound.WithMessage("db connection not found").WithStack(ctx)
}

func (pool *MongoDBPool) GetConfig(ctx context.Context) MongoDBPoolConfig {
	// return a config copy
	newCfg := MongoDBPoolConfig{}
	newCfg.Config = make([]*MongoConfig, len(pool.config.Config))
	copy(newCfg.Config, pool.config.Config)
	return newCfg
}

func (pool *MongoDBPool) String() string {
	sList := make([]string, 0)
	for _, c := range pool.config.Config {
		sList = append(sList, fmt.Sprintf("%s: %s", c.Name, c.MongoURI))
	}
	return fmt.Sprintf("MongoDBPool: \n%s", strings.Join(sList, ",\n"))
}
