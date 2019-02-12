package schemas

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/scottdelly/goql/models"
)

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
		"id": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(int)
		if ok {
			response := new(models.Song)
			err := json.Unmarshal(
				[]byte(fmt.Sprintf(`{"id": %d, "name": "tester", "artist_id": 2, "duration": 3}`, id)),
				response)
			return *response, err
		}
		return nil, nil
	},
}
