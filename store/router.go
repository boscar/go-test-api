package store

import (
	"log"
	"net/http"

	"go-test-api/config"

	"github.com/gorilla/mux"
)

// Imports
var controller Controller

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		"Authentication",
		"POST",
		"/get-token",
		controller.GetToken,
	},
	Route{
		"Index",
		"GET",
		"/",
		controller.Index,
	},
	Route{
		"AddProduct",
		"POST",
		"/AddProduct",
		AuthenticationMiddleware(controller.AddProduct),
	},
}

// More routes.....

// CreateRouter function configures a new router to the API
func CreateRouter(configuration config.Configuration) *mux.Router {
	controller = Controller{
		Repository: Repository{
			Config: configuration,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	initRoutes(router, &routes)

	return router
}

func initRoutes(router *mux.Router, routes *Routes) {
	for _, route := range *routes {
		log.Println(route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}
