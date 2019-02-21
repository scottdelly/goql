#### A graphql API written in Go to interface with a Postgres DB.

**Features:**

- Run in docker
- Packages are vendored
- Uses graphiql to auto-document query structures

**Understanding the Code:**

My hope is that much of this is self-explanatory, but here's a quick run down of the overall structure

Code is split into 4 layers:
- The Model (`models/`) layer defines db structures and conformance to interfaces
- The DB (`db_client/`)layer interfaces with postgres and presents a simple syntax for building queries
- The API (`api/`) layer interfaces with the http listener and connects the GraphQL layer to the DB layer
- The GraphQL Layer (`api/achemas/`) defines graphql schemas available to the API

The Model, DB, and GraphQL layers follow a similar pattern of one member of each layer per model type, plus some shared/helper member.

**How to Use:**
- `docker-compose up` will start the project
- Navigate in your browser to http://localhost:8080/
- Hit enter at the prompt to use `http://localhost:8080/graphql`
- Enter a query or browse the schema using the `docs` link on the top right
  
  Sample Query: 
```
query {
  user(id: 1) {
    id
    name
  }
  song(id: 42) {
    id
    name
    artist {
      id
      name
    }
  }
}
```

- The DB will bootstrap automatically, and be reset with reach restart of the server. To disable bootstrap, set "db_bootstrap" in config.json to false

See the mutation syntax to make API calls for:
- Creating Artists, Songs, and Users
- Liking artists and songs
  Sample Mutation: 
```
mutation {
  create_artist(name: "Dubfire") {
    success
    artist {
      id
      created
      name
    }
  }
}

```