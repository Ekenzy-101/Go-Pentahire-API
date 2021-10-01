package main

import (
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
	if !config.IsDevelopment {
		log.Printf("Server listening on port %v", config.Port)
	}

	helpers.ExitIfError(router.Run(":" + config.Port))
}
