package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/mnmtanish/go-graphiql"

	"github.com/scottdelly/goql/api/schemas"
)

type GQLApi struct {
	schema graphql.Schema
}

func (g *GQLApi) Start(host string) error {
	g.schema = g.startGQL()
	return g.startHttp(host)
}

func (g *GQLApi) Test() {
	// Query
	query := `
		{
			hello
		}
	`
	params := graphql.Params{Schema: g.schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}

func (g *GQLApi) startHttp(host string) error {
	http.HandleFunc("/gql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), g.schema)
		_ = json.NewEncoder(w).Encode(result)
	})
	http.HandleFunc("/", graphiql.ServeGraphiQL)

	fmt.Println(fmt.Sprintf("Server is running at %s", host))

	return http.ListenAndServe(host, nil)
}

func (g *GQLApi) startGQL() graphql.Schema {

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user":    schemas.UserQueryField,
				"song":    schemas.SongQueryField,
				"songs":   schemas.SongListField,
				"artist":  schemas.ArtistQueryField,
				"artists": schemas.ArtistListField,
			},
		},
	)

	var muatationType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutate",
			Fields: graphql.Fields{
				"like_artist": schemas.ArtistLikeMutation,
				"like_song":   schemas.SongLikeMutation,
			},
		},
	)

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: muatationType,
		},
	)

	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	return schema
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v\n", result.Errors)
	}
	return result
}
