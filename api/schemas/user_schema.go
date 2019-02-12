package schemas

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/models"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":   gqlIdField(),
			"name": gqlNameField(),
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var UserQueryField = &graphql.Field{
	Type:        userType,
	Description: "Get user by ID",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(int)
		if ok {
			response := new(models.User)
			err := json.Unmarshal(
				[]byte(fmt.Sprintf(`{"id": %d, "name": "tester", "email": "test@test.test"}`, id)),
				response)
			return *response, err
		}
		return nil, nil
	},
}
