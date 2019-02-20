package schemas

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
)

var DBC *db_client.DBClient

const IdField = "id"
const NameField = "name"
const CreatedField = "created"

func createGQLObject(name string, fields graphql.Fields) *graphql.Object {
	if fields == nil {
		fields = make(map[string]*graphql.Field)
	}

	fields[IdField] = gqlIdField()
	fields[NameField] = gqlNameField()
	fields[CreatedField] = gqlCreatedField()

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
	if id, ok := p.Args[IdField].(models.ModelId); ok {
		return id, nil
	}
	if id, ok := p.Source.(models.Identifiable); ok {
		return id.Identifier(), nil
	}
	return 0, errors.New(fmt.Sprintf("Failed to parse model Id from %+v", p.Source))
}
