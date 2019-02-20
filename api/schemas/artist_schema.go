package schemas

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/db_client"
	"github.com/scottdelly/goql/models"
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
	artistType.AddFieldConfig("liked_by", ArtistLikesField)
}

var ArtistQueryField = &graphql.Field{
	Type: artistType,
	Args: graphql.FieldConfigArgument{
		"id": modelIDArgConfig("Artist Id"),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if id, err := parseModelId(p); err != nil {
			return nil, err
		} else {
			return artistClient().GetArtistById(id)
		}
	},
}

var ArtistListField = &graphql.Field{
	Type:        graphql.NewList(artistType),
	Description: "List of Artists",
	Args: graphql.FieldConfigArgument{
		"limit": limitArgConfig(),
		"query": queryArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		query, ok := parseQuery(p)
		if ok {
			return artistClient().GetArtists(parseLimit(p), `"name" ilike $1`, fmt.Sprintf("%%%s%%", query))
		}
		return artistClient().GetArtists(parseLimit(p), nil)
	},
}

var ArtistLikesField = &graphql.Field{
	Type:        graphql.NewList(userType),
	Description: "Users who like this artist",
	Args: graphql.FieldConfigArgument{
		"limit": limitArgConfig(),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit := parseLimit(p)
		artist := p.Source.(*models.Artist)
		if userIds, err := likesClient().GetUsersByLikes(artist.Likes(), limit); err != nil {
			return nil, err
		} else if len(userIds) > 0 {
			users, err := userClient().GetUsers(limit, `id IN $1`, userIds)
			return users, err
		}
		return nil, nil
	},
}
