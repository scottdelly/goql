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

func buildLikeMessage(p graphql.ResolveParams, likeType models.LikeObjectType) (*models.LikeMessage, error) {
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

	var message = models.LikeMessage{User: user, ObjectType: likeType}

	switch likeType {
	case models.LikeTypeArtist:
		if artistId, ok := p.Args["artist_id"]; ok {
			if message.Object, err = artistClient().GetArtistById(artistId.(models.ModelId)); err != nil {
				return nil, err
			}
		}
	case models.LikeTypeSong:
		if songId, ok := p.Args["song_id"]; ok {
			if message.Object, err = songClient().GetSongById(songId.(models.ModelId)); err != nil {
				return nil, err
			}
		}
	}
	return &message, nil
}

var ArtistLikeQueryField = &graphql.Field{
	Type:        graphql.NewList(ModelIdScalar),
	Description: "Artists that the user likes",
	Args: graphql.FieldConfigArgument{
		"limit": limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		message, err := buildLikeMessage(p, models.LikeTypeArtist)
		if err != nil {
			return nil, err
		}
		var response []models.ModelId
		if response, err = likesClient().GetLikes(*message, limit); err != nil {
			return nil, err
		}
		return response, nil
	},
}

func likeMutationResponse(name string, fields graphql.Fields) *graphql.Object {
	if fields == nil {
		fields = make(map[string]*graphql.Field)
	}

	fields["success"] = &graphql.Field{
		Type: graphql.Boolean,
	}
	fields["user"] = &graphql.Field{
		Type: userType,
	}

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		},
	)
}

var ArtistLikeMutation = &graphql.Field{
	Type: likeMutationResponse("UserLikesArtist",
		graphql.Fields{"artist": &graphql.Field{
			Type: artistType,
		}},
	),
	Args: graphql.FieldConfigArgument{
		"artist_id": modelIDArgConfig("Artist Id"),
		"user_id":   modelIDArgConfig("User Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		message, err := buildLikeMessage(p, models.LikeTypeArtist)
		if err != nil {
			return nil, err
		}
		if err = likesClient().CreateUserLike(*message); err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true,
			"user":    message.User,
			"artist":  message.Object,
		}, nil
	},
}

var SongLikeQueryField = &graphql.Field{
	Type:        graphql.NewList(ModelIdScalar),
	Description: "Songs that the user likes",
	Args: graphql.FieldConfigArgument{
		"limit": limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		message, err := buildLikeMessage(p, models.LikeTypeSong)
		if err != nil {
			return nil, err
		}
		var response []models.ModelId
		if response, err = likesClient().GetLikes(*message, limit); err != nil {
			return nil, err
		}
		return response, nil
	},
}

var SongLikeMutation = &graphql.Field{
	Type: likeMutationResponse("UserLikesSong",
		graphql.Fields{"song": &graphql.Field{
			Type: songType,
		}},
	),
	Args: graphql.FieldConfigArgument{
		"song_id": modelIDArgConfig("Song Id"),
		"user_id": modelIDArgConfig("User Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		message, err := buildLikeMessage(p, models.LikeTypeSong)
		if err != nil {
			return nil, err
		}
		if err = likesClient().CreateUserLike(*message); err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true,
			"user":    message.User,
			"song":    message.Object,
		}, nil
	},
}
