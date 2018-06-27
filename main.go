package main

import (
	"log"
	"net/http"

	"github.com/boscar/go-test-api/config"
	"github.com/boscar/go-test-api/store"

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
	err := gonfig.GetConf("config/config.atlas.development.json", &configuration)
	if err != nil {
		log.Fatal(err)
	}
	return configuration
}
