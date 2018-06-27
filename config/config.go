package config

//Configuration ...
type Configuration struct {
	Port             string   `json:"Port"`
	ConnectionString string   `json:"ConnectionString"`
	DatabaseName     string   `json:"DatabaseName"`
	Hosts            []string `json:"Hosts"`
	Password         string   `json:"Password"`
	DBUser           string   `json:"DBUser"`
}
