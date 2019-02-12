package schemas

import (
	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
)

func userClient() *db_client.UserClient {
	return db_client.NewUserClient(DBC)
}

var userType = createGQLObject("User",
	graphql.Fields{
		"email": &graphql.Field{
			Type: graphql.String,
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
		return userClient().GetUserById(parseModelId(p))
	},
}
