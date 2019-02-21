package schemas

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/models"
)

var artistType = createGQLObject("Artist",
	graphql.Fields{
		"like_count": gqlLikeCountField(),
	},
)

func init() {
	artistType.AddFieldConfig("songs", SongListField)
	artistType.AddFieldConfig("liked_by", ArtistLikesField)
}

var ArtistQueryField = &graphql.Field{
	Type: artistType,
	Args: graphql.FieldConfigArgument{
		IdField: modelIDArgConfig("Artist Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if id, err := parseModelId(p); err != nil {
			return nil, err
		} else {
			return ArtistClient.GetArtistById(id)
		}
	},
}

var ArtistListField = &graphql.Field{
	Type:        graphql.NewList(artistType),
	Description: "List of Artists",
	Args: graphql.FieldConfigArgument{
		LimitArg: limitArgConfig(),
		QueryArg: queryArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		query, ok := parseQuery(p)
		if ok {
			return ArtistClient.GetArtists(parseLimit(p), `"name" ilike $1`, fmt.Sprintf("%%%s%%", query))
		}
		return ArtistClient.GetArtists(parseLimit(p), nil)
	},
}

var ArtistLikesField = &graphql.Field{
	Type:        graphql.NewList(userType),
	Description: "Users who like this artist",
	Args: graphql.FieldConfigArgument{
		LimitArg: limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return parseUsersWhoLike(p)
	},
}

var ArtistCreateMutation = &graphql.Field{
	Type: mutationResponse("CreateArtist",
		graphql.Fields{
			"artist": &graphql.Field{
				Type: artistType,
			},
		},
	),
	Args: graphql.FieldConfigArgument{
		NameField: &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		artist := new(models.Artist)
		artist.Name = p.Args[NameField].(string)
		if err := ArtistClient.Create(artist); err != nil {
			return nil, graphql.NewLocatedError(err, nil)
		}
		return map[string]interface{}{
			"success": true,
			"artist":  artist,
		}, nil
	},
}
