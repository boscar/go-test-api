package store

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"go-test-api/config"

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
	session, err := connectToServer(r.Config)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err.Error())
		return nil
	}

	defer session.Close()

	collection := session.DB(r.Config.DatabaseName).C(STORE_COLLECTION)
	results := Products{}

	err = collection.Find(nil).All(&results)
	if err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// AddProduct adds a Product in the DB
func (r Repository) AddProduct(product Product) bool {
	session, err := connectToServer(r.Config)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
		return false
	}

	defer session.Close()

	productID++
	product.ID = productID

	err = session.DB(r.Config.DatabaseName).C(STORE_COLLECTION).Insert(product)
	if err != nil {
		fmt.Println("Failed to insert product:", err)
		return false
	}

	fmt.Println("Added New Product ID- ", product.ID)

	return true
}

func connectToServer(config config.Configuration) (*mgo.Session, error) {
	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs:    config.Hosts,
		Timeout:  30 * time.Second,
		Username: config.DBUser,
		Password: config.Password,
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	return mgo.DialWithInfo(dialInfo)
}
