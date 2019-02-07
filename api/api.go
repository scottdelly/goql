package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"

	"scottdelly/goql/models"
)

type GQLApi struct {
	schema graphql.Schema
}

func (g *GQLApi) Start(host string) {
	g.schema = g.startGQL()
	g.startHttp(host)
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

func (g *GQLApi) startHttp(host string) {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), g.schema)
		_ = json.NewEncoder(w).Encode(result)
	})
	fmt.Println(fmt.Sprintf("Server is running at %s", host))

	err := http.ListenAndServe(host, nil)
	if err != nil {
		panic(err)
	}
}

func (g *GQLApi) startGQL() graphql.Schema {

	var userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: userType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						response := new(models.User)
						err := json.Unmarshal([]byte(`{"id":"tester", "name":"tester"}`), response)
						println(fmt.Sprintf("%+v", response))
						return response, err
					},
				},
			},
		})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
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
