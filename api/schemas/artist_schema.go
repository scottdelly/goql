package schemas

import (
	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
)

func artistClient() *db_client.ArtistClient {
	return db_client.NewArtistClient(DBC)
}

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
		"id": modelIDArgumentConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(models.ModelId)
		if ok {
			artist, err := artistClient().GetArtistById(id)
			return *artist, err
		}
		return nil, nil
	},
}
