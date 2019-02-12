package schemas

import (
	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
)

func songClient() *db_client.SongClient {
	return db_client.NewSongClient(DBC)
}

var songType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Song",
		Fields: graphql.Fields{
			"id":   gqlIdField(),
			"name": gqlNameField(),
			"artist_id": &graphql.Field{
				Type: ModelIdScalar,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(models.Song).ArtistId, nil
				},
			},
			"duration": &graphql.Field{
				Type: DurationScalar,
			},
		},
	},
)

var SongQueryField = &graphql.Field{
	Type: songType,
	Args: graphql.FieldConfigArgument{
		"id": modelIDArgumentConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(models.ModelId)
		if ok {
			return songClient().GetSongById(id)
		}
		return nil, nil
	},
}
