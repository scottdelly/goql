package main

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
	"log"

	"github.com/scottdelly/goql/api"
	"github.com/scottdelly/goql/db_client"

	"github.com/spf13/pflag"
)

func main() {
	log.Println("==> Starting up")
	viper.SetDefault("api_host_name", ":8080")

	// using standard library "flag" package
	flag.String("db_host_name", "localhost", "override host name for database")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Errorf("fatal error binding to cmd line %s", err))
	}
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	apiHostName := viper.GetString("api_host_name")

	apiErr := make(chan error)
	go func() {
		gqlAPI := new(api.GQLApi)
		apiErr <- gqlAPI.Start(apiHostName)
	}()

	dbHostName := viper.GetString("db_host_name")
	dbUser := viper.GetString("db_user")
	dbPass := viper.GetString("db_pass")
	dbBootstrap := viper.GetBool("db_bootstrap")

	go func() {
		dbc := new(db_client.DBClient)
		api.SetDbClient(dbc)
		dbc.Start(dbUser, dbPass, dbHostName, dbBootstrap)
	}()

	err = <-apiErr
	log.Fatal(err.Error())
}
