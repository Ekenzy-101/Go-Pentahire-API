package main

import (
	"fmt"
	"log"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/routes"
	"github.com/Ekenzy-101/Pentahire-API/services"
)

func main() {
	pool := services.CreatePostgresConnectionPool()
	defer pool.Close()

	services.CreateRedisClient()
	router := routes.SetupRouter()
	host := "127.0.0.1"
	if config.IsProduction {
		log.Printf("Server listening on port %v", config.Port)
		host = ""
	}

	err := router.Run(fmt.Sprintf("%v:%v", host, config.Port))
	helpers.ExitIfError(err)
}
