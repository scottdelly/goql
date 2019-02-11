package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/scottdelly/goql/api"
	"github.com/scottdelly/goql/db_client"
)

func main() {
	viper.SetDefault("host", ":8000")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	host := viper.GetString("host")

	gqlAPI := new(api.GQLApi)
	gqlAPI.Start(host)
	db := new(db_client.DBClient)
	db.Start()
}
