Welcome to my Golang graphQL demo.

This repo will contain a working application that creates a graphQL interface to a SQL database.

Progress updates:

commit ec74244c:
- Pulls in a docker image for postgres and the entire app can be started with 'docker-compose up'
- Connects db layer to actual db

Commit f4fd493b:
- Update the API to resolve graphql queries with all current types
- adds graphiql server for auto-documentation

Commit 1996e55:
- Update dbclient to have a CRUD interface for easily working with models

Initial commit:
- GraphQL skeleton in place
- Bare-bones models in place
- DB initializer in place but untested