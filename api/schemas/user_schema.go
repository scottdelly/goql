package schemas

import (
	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
)

func userClient() *db_client.UserClient {
	return db_client.NewUserClient(DBC)
}

func likesClient() *db_client.LikesClient {
	return db_client.NewLikesClient(DBC)
}

var userType = createGQLObject("User",
	graphql.Fields{
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"liked_artists": ArtistLikeQueryField,
		"liked_songs":   SongLikeQueryField,
	},
)

var UserQueryField = &graphql.Field{
	Type:        userType,
	Description: "Get user by Id",
	Args: graphql.FieldConfigArgument{
		"id": modelIDArgConfig("User Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if id, err := parseModelId(p); err != nil {
			return nil, err
		} else {
			return userClient().GetUserById(id)
		}
	},
}

func userFromParams(p graphql.ResolveParams) (*models.User, error) {
	var err error
	var userId models.ModelId
	if userArg, ok := p.Args["user_id"]; ok {
		userId = userArg.(models.ModelId)
	} else {
		userId = p.Source.(*models.User).Id
	}
	var user *models.User
	if user, err = userClient().GetUserById(userId); err != nil {
		return nil, err
	}
	return user, nil
}

var ArtistLikeQueryField = &graphql.Field{
	Type:        graphql.NewList(artistType),
	Description: "Artists that the user likes",
	Args: graphql.FieldConfigArgument{
		"limit": limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		user, err := userFromParams(p)
		if err != nil {
			return nil, err
		}
		var response []models.ModelId
		if response, err = likesClient().GetLikesForUser(user.LikesOfType(models.LikeTypeArtist), limit); err != nil {
			return nil, err
		}
		if len(response) > 0 {
			var artists []*models.Artist
			artists, err = artistClient().GetArtists(limit, `"id" IN $1`, response)
			if err != nil {
				return nil, err
			}
			return artists, nil
		}
		return nil, nil
	},
}

var ArtistLikeMutation = &graphql.Field{
	Type: mutationResponse("UserLikesArtist",
		graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
			},
			"artist": &graphql.Field{
				Type: artistType,
			},
		},
	),
	Args: graphql.FieldConfigArgument{
		"artist_id": modelIDArgConfig("Artist Id"),
		"user_id":   modelIDArgConfig("User Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		user, err := userFromParams(p)
		if err != nil {
			return nil, err
		}
		var artist *models.Artist
		if artistId, ok := p.Args["artist_id"]; ok {
			if artist, err = artistClient().GetArtistById(artistId.(models.ModelId)); err != nil {
				return nil, err
			}
		}
		if err = likesClient().CreateUserLike(user.LikeArtist(artist)); err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true,
			"user":    user,
			"artist":  artist,
		}, nil
	},
}

var SongLikeQueryField = &graphql.Field{
	Type:        graphql.NewList(songType),
	Description: "Songs that the user likes",
	Args: graphql.FieldConfigArgument{
		"limit": limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		user, err := userFromParams(p)
		if err != nil {
			return nil, err
		}
		var response []models.ModelId
		if response, err = likesClient().GetLikesForUser(user.LikesOfType(models.LikeTypeSong), limit); err != nil {
			return nil, err
		}
		if len(response) > 0 {
			var songs []*models.Song
			songs, err = songClient().GetSongs(limit, `"id" IN $1`, response)
			if err != nil {
				return nil, err
			}
			return songs, nil
		}
		return response, nil
	},
}

var SongLikeMutation = &graphql.Field{
	Type: mutationResponse("UserLikesSong",
		graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
			},
			"song": &graphql.Field{
				Type: songType,
			},
		},
	),
	Args: graphql.FieldConfigArgument{
		"song_id": modelIDArgConfig("Song Id"),
		"user_id": modelIDArgConfig("User Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		user, err := userFromParams(p)
		if err != nil {
			return nil, err
		}
		var song *models.Song
		if songId, ok := p.Args["song_id"]; ok {
			if song, err = songClient().GetSongById(songId.(models.ModelId)); err != nil {
				return nil, err
			}
		}
		if err = likesClient().CreateUserLike(user.LikeSong(song)); err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true,
			"user":    user,
			"song":    song,
		}, nil
	},
}
