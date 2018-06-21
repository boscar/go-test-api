package store

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "mongodb://localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "dummyStore"

// COLLECTION is the name of the collection in DB
const COLLECTION = "store"

var productID = 10

// GetProducts returns the list of Products
func (r Repository) GetProducts() Products {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
		return nil
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Products{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// AddProduct adds a Product in the DB
func (r Repository) AddProduct(product Product) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	productID++
	product.ID = productID
	session.DB(DBNAME).C(COLLECTION).Insert(product)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Product ID- ", product.ID)

	return true
}
