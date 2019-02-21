package schemas

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/models"
)

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
		"like_count": gqlLikeCountField(),
	},
)

func init() {
	songType.AddFieldConfig("artist", &graphql.Field{
		Type: artistType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			artistId := p.Source.(*models.Song).ArtistId
			return ArtistClient.GetArtistById(artistId)
		},
	})
	songType.AddFieldConfig("liked_by", SongLikesField)
}

var SongQueryField = &graphql.Field{
	Type: songType,
	Args: graphql.FieldConfigArgument{
		IdField: modelIDArgConfig("Song Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if id, err := parseModelId(p); err != nil {
			return nil, err
		} else {
			return SongClient.GetSongById(id)
		}
	},
}

var SongListField = &graphql.Field{
	Type:        graphql.NewList(songType),
	Description: "List of Songs",
	Args: graphql.FieldConfigArgument{
		LimitArg: limitArgConfig(),
		QueryArg: queryArgConfig(),
		"artist_id": &graphql.ArgumentConfig{
			Type: ModelIdScalar,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if query, ok := parseQuery(p); ok {
			return SongClient.GetSongs(parseLimit(p), `"name" ilike $1`, fmt.Sprintf("%%%s%%", query))
		}
		artistId := models.ModelId(-1)
		if artist, ok := p.Source.(*models.Artist); ok {
			artistId = artist.Id
		} else if id, ok := p.Args["artist_id"].(models.ModelId); ok {
			artistId = id
		}
		if artistId > -1 {
			return SongClient.GetSongs(parseLimit(p), `"artist_id" = $1`, artistId)
		} else {
			return SongClient.GetSongs(parseLimit(p), nil)
		}
	},
}

var SongLikesField = &graphql.Field{
	Type:        graphql.NewList(userType),
	Description: "Users who like this song",
	Args: graphql.FieldConfigArgument{
		LimitArg: limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return parseUsersWhoLike(p)
	},
}

var SongCreateMutation = &graphql.Field{
	Type: mutationResponse("CreateSong",
		graphql.Fields{
			"song": &graphql.Field{
				Type: songType,
			},
		},
	),
	Args: graphql.FieldConfigArgument{
		NameField: &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"artist_id": modelIDArgConfig("Artist Id"),
		"duration": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(DurationScalar),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		song := new(models.Song)
		song.Name = p.Args[NameField].(string)
		song.ArtistId = p.Args["artist_id"].(models.ModelId)
		song.Duration = p.Args["duration"].(time.Duration)
		if err := SongClient.Create(song); err != nil {
			return nil, graphql.NewLocatedError(err, nil)
		}
		return map[string]interface{}{
			"success": true,
			"song":    song,
		}, nil
	},
}
