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
		IdField: modelIDArgConfig("User Id"),
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
		LimitArg: limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		response, err := parseLikesForUser(p, models.LikeTypeArtist)
		if err != nil || len(response) == 0 {
			return nil, err
		}
		artists, err := artistClient().GetArtists(limit, `"id" IN $1`, response)
		return artists, err
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
		return createLike(p, models.LikeTypeArtist)
	},
}

var SongLikeQueryField = &graphql.Field{
	Type:        graphql.NewList(songType),
	Description: "Songs that the user likes",
	Args: graphql.FieldConfigArgument{
		LimitArg: limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		response, err := parseLikesForUser(p, models.LikeTypeSong)
		if err != nil || len(response) == 0 {
			return nil, err
		}
		songs, err := songClient().GetSongs(limit, `"id" IN $1`, response)
		return songs, err
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
		return createLike(p, models.LikeTypeSong)
	},
}
