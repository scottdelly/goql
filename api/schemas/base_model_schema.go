package schemas

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
)

var DBC *db_client.DBClient

func createGQLObject(name string, fields graphql.Fields) *graphql.Object {

	if fields == nil {
		fields = make(map[string]*graphql.Field)
	}

	fields["id"] = gqlIdField()
	fields["name"] = gqlNameField()
	fields["created"] = gqlCreatedField()

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		},
	)
}

func gqlIdField() *graphql.Field {
	return &graphql.Field{
		Type: ModelIdScalar,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return p.Source.(models.Identifiable).Identifier(), nil
		},
	}
}

func gqlNameField() *graphql.Field {
	return &graphql.Field{
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return p.Source.(models.Nameable).NameValue(), nil
		},
	}
}

func gqlCreatedField() *graphql.Field {
	return &graphql.Field{
		Type: graphql.DateTime,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return p.Source.(models.Historical).DateValue(), nil
		},
		Description: "Date created",
	}
}

func modelIDArgConfig(description string) *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type:        ModelIdScalar,
		Description: description,
	}
}

func parseModelId(p graphql.ResolveParams) (models.ModelId, error) {
	if id, ok := p.Args["id"].(models.ModelId); ok {
		return id, nil
	}
	if id, ok := p.Source.(models.Identifiable); ok {
		return id.Identifier(), nil
	}
	return 0, errors.New(fmt.Sprintf("Failed to parse model Id from %+v", p.Source))
}

func limitArgConfig() *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 10,
		Description:  "Maximum number of results returned",
	}
}

func parseLimit(p graphql.ResolveParams) uint64 {
	return uint64(p.Args["limit"].(int))
}

func queryArgConfig() *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "Search for query string",
	}
}

func parseQuery(p graphql.ResolveParams) (string, bool) {
	val, ok := p.Args["query"].(string)
	return val, ok
}
