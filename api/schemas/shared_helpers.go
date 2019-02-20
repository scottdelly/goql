package schemas

import (
	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/models"
)

const SuccessField = "success"

const LimitArg = "limit"
const QueryArg = "query"

//Limits
func limitArgConfig() *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 10,
		Description:  "Maximum number of results returned",
	}
}

func parseLimit(p graphql.ResolveParams) uint64 {
	return uint64(p.Args[LimitArg].(int))
}

//Queries
func queryArgConfig() *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "Search for query string",
	}
}

func parseQuery(p graphql.ResolveParams) (string, bool) {
	val, ok := p.Args[QueryArg].(string)
	return val, ok
}

//Mutation
func mutationResponse(name string, fields graphql.Fields) *graphql.Object {
	if fields == nil {
		fields = make(map[string]*graphql.Field)
	}

	fields[SuccessField] = &graphql.Field{
		Type: graphql.Boolean,
	}

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		},
	)
}

//Likes
func gqlLikeCountField() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Int,
		Description: "Number of likes on the recipient",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			likable := p.Source.(models.Likable)
			count, err := likesClient().LikeCount(models.LikesOn(likable))
			return count, err
		},
	}
}

func parseLikesForUser(p graphql.ResolveParams, t models.LikeObjectType) ([]models.ModelId, error) {
	user, err := userFromParams(p)
	limit := parseLimit(p)
	if err != nil {
		return nil, err
	}
	var response []models.ModelId
	if response, err = likesClient().GetLikesForUser(user.LikesOfType(t), limit); err != nil {
		return nil, err
	}
	return response, nil
}

func parseUsersWhoLike(p graphql.ResolveParams) ([]*models.User, error) {
	limit := parseLimit(p)
	likable := p.Source.(models.Likable)
	if userIds, err := likesClient().GetUsersByLikes(models.LikesOn(likable), limit); err != nil {
		return nil, err
	} else if len(userIds) > 0 {
		users, err := userClient().GetUsers(limit, `id IN $1`, userIds)
		return users, err
	}
	return nil, nil
}

func createLike(p graphql.ResolveParams, t models.LikeObjectType) (map[string]interface{}, error) {
	user, err := userFromParams(p)
	if err != nil {
		return nil, err
	}
	var objectId models.ModelId
	var likeObjectKey string
	switch t {
	case models.LikeTypeArtist:
		likeObjectKey = "artist"
		objectId = p.Args["artist_id"].(models.ModelId)
	case models.LikeTypeSong:
		likeObjectKey = "song"
		objectId = p.Args["song_id"].(models.ModelId)

	}
	var likable models.Likable
	if likable, err = likesClient().ResolveObject(objectId, t); err != nil {
		return nil, err
	}

	if err = likesClient().CreateUserLike(user.LikeObject(likable)); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"success":     true,
		"user":        user,
		likeObjectKey: likable,
	}, nil
}
