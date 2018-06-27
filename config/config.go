package config

//Configuration ...
type Configuration struct {
	Port             string
	ConnectionString string
	DatabaseName     string
	UserName         string `json:"UserName"`
	Password         string `json:"Password"`
	DBUser           string `json:"DBUser"`
}
