package services

import (
	"context"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/env"
)

type DatabaseService struct {
	*ent.Client
}

func NewDatabaseService(driver string, dataSourceName string) (*DatabaseService, error) {
	client, err := ent.Open(env.DatabaseDriver, env.DatabaseEndpoint)
	if err != nil {
		return nil, err
	}

	return &DatabaseService{
		Client: client,
	}, nil
}

func (s *DatabaseService) Close() error {
	return s.Client.Close()
}

func (s *DatabaseService) AutoMigrate() error {
	return s.Client.Schema.Create(context.Background())
}
