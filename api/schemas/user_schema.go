package schemas

import (
	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
)

func userClient() *db_client.UserClient {
	return db_client.NewUserClient(DBC)
}

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
		"id": modelIDArgumentConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(models.ModelId)
		if ok {
			return userClient().GetUserById(id)
		}
		return nil, nil
	},
}
