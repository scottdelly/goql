package schemas

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/models"
)

var artistType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Artist",
		Fields: graphql.Fields{
			"id":   gqlIdField(),
			"name": gqlNameField(),
			"like_count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var ArtistQueryField = &graphql.Field{
	Type: artistType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(int)
		if ok {
			response := new(models.Artist)
			err := json.Unmarshal(
				[]byte(fmt.Sprintf(`{"id": %d,  "name": "tester", "like_count": 4}`, id)),
				response)
			return *response, err
		}
		return nil, nil
	},
}
