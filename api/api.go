package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/mnmtanish/go-graphiql"

	"github.com/scottdelly/goql/api/schemas"
	"github.com/scottdelly/goql/db_client"
)

type GQLApi struct {
	schema graphql.Schema
}

var dBC *db_client.DBClient

func (g *GQLApi) Start(host string) error {
	log.Println("==> API Starting")
	g.schema = g.startGQL()
	return g.startHttp(host)
}

func (g *GQLApi) startHttp(host string) error {

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var query string
		switch strings.ToLower(r.Method) {
		case "get":
			query = r.URL.Query().Get("query")
		case "post":
			decoder := json.NewDecoder(r.Body)
			var graphQL map[string]interface{}
			if err := decoder.Decode(&graphQL); err != nil {
				panic(err)
			}
			query = string(graphQL["query"].(string))
		}
		result := executeQuery(query, g.schema)
		_ = json.NewEncoder(w).Encode(result)
	})
	http.HandleFunc("/", graphiql.ServeGraphiQL)

	log.Printf("==> API listening at %s", host)

	return http.ListenAndServe(host, nil)
}

func SetDbClient(client *db_client.DBClient) {
	dBC = client
	schemas.ArtistClient = db_client.NewArtistClient(client)
	schemas.SongClient = db_client.NewSongClient(client)
	schemas.UserClient = db_client.NewUserClient(client)
	schemas.LikesClient = db_client.NewLikesClient(client)
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

	var mutationType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutate",
			Fields: graphql.Fields{
				"create_artist": schemas.ArtistCreateMutation,
				"create_song":   schemas.SongCreateMutation,
				"create_user":   schemas.UserCreateMutation,
				"like_artist":   schemas.ArtistLikeMutation,
				"like_song":     schemas.SongLikeMutation,
			},
		},
	)

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
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
		log.Printf("wrong result, unexpected errors: %v\n", result.Errors)
	}
	return result
}
