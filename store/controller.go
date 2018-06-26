package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

//Controller ...
type Controller struct {
	Repository Repository
}

/*AuthenticationMiddleware ...
Middleware handler to handle all requests for authentication */
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		authorizationHeader := request.Header.Get("authorization")

		if authorizationHeader == "" {
			json.NewEncoder(responseWriter).Encode(Exception{Message: "An authorization header is required."})
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")

		if len(bearerToken) != 2 {
			json.NewEncoder(responseWriter).Encode(Exception{Message: "Invalid authorization header format."})
			return
		}

		claims, error := getTokenClaims(&bearerToken[1])

		if error == nil {
			log.Println("TOKEN WAS VALID")
			context.Set(request, "decoded", claims)
			next(responseWriter, request)
		} else {
			json.NewEncoder(responseWriter).Encode(Exception{Message: error.Error()})
		}
	})
}

func getTokenClaims(tokenString *string) (jwt.Claims, error) {
	token, error := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if error != nil {
		return nil, error
	}
	if token.Valid {
		return token.Claims, nil
	}
	error = errors.New("Invalid authorization token")
	return nil, error
}

//GetToken ... Get Authentication token GET /
func (c *Controller) GetToken(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	log.Println("Username: " + user.Username)
	log.Println("Password: " + user.Password)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

// Index GET /
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	products := c.Repository.GetProducts()
	data, _ := json.Marshal(products)
	setJSONHeader(w)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddProduct POST /
func (c *Controller) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request

	log.Println(body)

	if err != nil {
		log.Fatalln("Error AddProduct", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddProduct", err)
	}

	if err := json.Unmarshal(body, &product); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	log.Println(product)
	success := c.Repository.AddProduct(product) // adds the product to the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setJSONHeader(w)
	w.WriteHeader(http.StatusCreated)
	return
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}
