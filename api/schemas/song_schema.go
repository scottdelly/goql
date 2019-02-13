package schemas

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
)

func songClient() *db_client.SongClient {
	return db_client.NewSongClient(DBC)
}

var songType = createGQLObject("Song",
	graphql.Fields{
		"artist_id": &graphql.Field{
			Type: ModelIdScalar,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(*models.Song).ArtistId, nil
			},
		},
		"duration": &graphql.Field{
			Type: DurationScalar,
		},
	},
)

func init() {
	songType.AddFieldConfig("artist", &graphql.Field{
		Type: artistType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			artistId := p.Source.(*models.Song).ArtistId
			return artistClient().GetArtistById(artistId)
		},
	})
}

var SongQueryField = &graphql.Field{
	Type: songType,
	Args: graphql.FieldConfigArgument{
		"id": modelIDArgumentConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if id, err := parseModelId(p); err != nil {
			return nil, err
		} else {
			return songClient().GetSongById(id)
		}
	},
}

var SongListField = &graphql.Field{
	Type:        graphql.NewList(songType),
	Description: "List of Songs",
	Args: graphql.FieldConfigArgument{
		"limit": limitFieldConfig(),
		"query": queryFieldConfig(),
		"artist_id": &graphql.ArgumentConfig{
			Type: ModelIdScalar,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if query, ok := parseQuery(p); ok {
			return songClient().GetSongs(parseLimit(p), `"name" ilike $1`, fmt.Sprintf("%%%s%%", query))
		}
		artistId := models.ModelId(-1)
		if artist, ok := p.Source.(*models.Artist); ok {
			artistId = artist.ID
		} else if id, ok := p.Args["artist_id"].(models.ModelId); ok {
			artistId = id
		}
		if artistId > -1 {
			return songClient().GetSongs(parseLimit(p), `"artist_id" = $1`, artistId)
		} else {
			return songClient().GetSongs(parseLimit(p), nil)
		}
	},
}
