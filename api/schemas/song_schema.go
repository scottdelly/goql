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
	songType.AddFieldConfig("liked_by", SongLikesField)
}

var SongQueryField = &graphql.Field{
	Type: songType,
	Args: graphql.FieldConfigArgument{
		"id": modelIDArgConfig("Song Id"),
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
		"limit": limitArgConfig(),
		"query": queryArgConfig(),
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
			artistId = artist.Id
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

var SongLikesField = &graphql.Field{
	Type:        graphql.NewList(userType),
	Description: "Users who like this song",
	Args: graphql.FieldConfigArgument{
		"limit": limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		song := p.Source.(*models.Song)
		if userIds, err := likesClient().GetUsersByLikes(song.Likes(), limit); err != nil {
			return nil, err
		} else if len(userIds) > 0 {
			users, err := userClient().GetUsers(limit, `id IN $1`, userIds)
			return users, err
		}
		return nil, nil
	},
}
