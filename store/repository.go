package store

import (
	"fmt"
	"go-test-api/config"
	"log"

	mgo "gopkg.in/mgo.v2"
)

//Repository ...
type Repository struct {
	Config config.Configuration
}

// STORE_COLLECTION is the name of the collection in DB
const STORE_COLLECTION = "store"

var productID = 10

// GetProducts returns the list of Products
func (r Repository) GetProducts() Products {
	session, err := mgo.Dial(r.Config.ConnectionString)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
		return nil
	}

	defer session.Close()

	c := session.DB(r.Config.DatabaseName).C(STORE_COLLECTION)
	results := Products{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// AddProduct adds a Product in the DB
func (r Repository) AddProduct(product Product) bool {
	session, err := mgo.Dial(r.Config.ConnectionString)
	defer session.Close()

	productID++
	product.ID = productID
	session.DB(r.Config.DatabaseName).C(STORE_COLLECTION).Insert(product)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Product ID- ", product.ID)

	return true
}
