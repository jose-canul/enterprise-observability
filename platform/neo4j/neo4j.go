package neo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func InsertResource(driver neo4j.Driver, rss Resource) (*Resource, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return createResource(tx, rss)
	})
	if err != nil {
		return nil, err
	}
	return result.(*Resource), nil
}

func createResource(tx neo4j.Transaction, rss Resource) (interface{}, error) {
	resource := map[string]interface{}{
		"id":   rss.Id,
		"name": rss.Name,
	}

	query := "CREATE (n:Resource { id: $id, name: $name }) RETURN n.id, n.name"

	records, err := tx.Run(query, resource)
	if err != nil {
		return nil, err
	}

	record, err := records.Single()
	if err != nil {
		return nil, err
	}

	return &Resource{
		Id:   record.Values[0].(string),
		Name: record.Values[1].(string),
	}, nil
}

func QueryResource(driver neo4j.Driver, rss Resource) (*Resource, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return fetchResource(tx, rss)
	})
	if err != nil {
		return nil, err
	}
	return result.(*Resource), nil
}

func fetchResource(tx neo4j.Transaction, rss Resource) (interface{}, error) {
	resource := map[string]interface{}{
		"id":   rss.Id,
		"name": rss.Name,
	}

	query := "MATCH (n:Resource { id: $id, name: $name }) RETURN n.id, n.name"

	records, err := tx.Run(query, resource)
	if err != nil {
		return nil, err
	}

	record, err := records.Single()
	if err != nil {
		return nil, err
	}

	return &Resource{
		Id:   record.Values[0].(string),
		Name: record.Values[1].(string),
	}, nil
}

func InsertRelationship(driver neo4j.Driver, rss1 Resource, rel Relationship, rss2 Resource) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return nil, createRelationship(tx, rss1, rel, rss2)
	})
	if err != nil {
		return err
	}
	return nil
}

func createRelationship(tx neo4j.Transaction, rss1 Resource, rel Relationship, rss2 Resource) error {
	relationship := map[string]interface{}{
		"id1": rss1.Id,
		"id2": rss2.Id,
	}

	query := fmt.Sprintf("MATCH (r1:Resource {id: $id1}), (r2:Resource {id: $id2}) MERGE (r1)-[r:%v]->(r2)", rel.Name)
	_, err := tx.Run(query, relationship)
	if err != nil {
		return err
	}
	return nil
}

type Resource struct {
	Id   string
	Name string
}

type Relationship struct {
	Name string
}
