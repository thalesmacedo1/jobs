package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jClient struct {
	Driver neo4j.DriverWithContext
}

func NewNeo4jClient(uri, username, password string) (*Neo4jClient, error) {

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar o driver Neo4j: %w", err)
	}

	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		return nil, fmt.Errorf("falha ao verificar a conectividade Neo4j: %w", err)
	}

	return &Neo4jClient{
		Driver: driver,
	}, nil
}

func (c *Neo4jClient) Close() error {
	return c.Driver.Close(context.Background())
}
