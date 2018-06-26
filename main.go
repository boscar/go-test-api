package main

import (
	"github.com/boscar/go-test-api/config"
	"github.com/boscar/go-test-api/store"

	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/tkanos/gonfig"
)

func main() {
	configuration := setupConfiguration()

	port := configuration.Port
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := store.CreateRouter(configuration)

	// These two lines are important if you're designing a front-end to utilise this API methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	// Launch server with CORS validations
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods)(router)))
}

func setupConfiguration() config.Configuration {
	configuration := config.Configuration{}
	err := gonfig.GetConf("config/config.docker.json", &configuration)
	if err != nil {
		log.Fatal(err)
	}
	return configuration
}
