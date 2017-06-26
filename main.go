package main

import (
	"github.com/spf13/viper"
	"log"
	"amigo-tech-test/service"
	"amigo-tech-test/util"
)

func main() {
	viper.SetConfigFile("./config/conf.json")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	a := App{}
	a.Initialise(
		&service.MessageServiceRouter{},
		util.PostgresDatabaseConnector{},
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.database"))
	a.Run(util.DefaultHttpServer{}, viper.GetString("app.port"))
}