package main

import (
	"example/vigovia-itenary-api/routes"
	"example/vigovia-itenary-api/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main(){

	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println(".env file failed to load. Using default environment variables.")
	}

	//initializes the gin router
	router:=gin.Default()

	//loads the configuration
	cfg:=config.NewConfig()

	//calls the SetupRoutes to set up the routes and pass config information
	routes.SetupRoutes(router,cfg)
	
	//starts the HTTP server 
	router.Run(cfg.ServerAddress)
}

