package main

import (
	auth "authBack"
	"authBack/pkg/handler"
	"log"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error while initialising config: %s", err.Error())
	}

	

	srv := new(auth.Server)
	log.Println(viper.GetString("port"))
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatal("Error while starting server: ", err)
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
