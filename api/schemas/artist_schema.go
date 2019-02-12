package schemas

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
)

func artistClient() *db_client.ArtistClient {
	return db_client.NewArtistClient(DBC)
}

var artistType = createGQLObject("Artist",
	graphql.Fields{
		"like_count": &graphql.Field{
			Type: graphql.Int,
		},
	},
)

func init() {
	artistType.AddFieldConfig("songs", SongListField)
}

var ArtistQueryField = &graphql.Field{
	Type: artistType,
	Args: graphql.FieldConfigArgument{
		"id": modelIDArgumentConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return artistClient().GetArtistById(parseModelId(p))
	},
}

var ArtistListField = &graphql.Field{
	Type:        graphql.NewList(songType),
	Description: "List of Artists",
	Args: graphql.FieldConfigArgument{
		"limit": limitFieldConfig(),
		"query": queryFieldConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		query, ok := parseQuery(p)
		if ok {
			return artistClient().GetArtists(parseLimit(p), `"name" ilike $1`, fmt.Sprintf("%%%s%%", query))
		}
		return artistClient().GetArtists(parseLimit(p), nil)
	},
}
