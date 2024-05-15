package dns

import (
	"encoding/json"
	"fmt"
	"os"
)

type dnsentry struct {
	Dns string `json:"dns"`
	Ip  string `json:"ip"`
}

type Alldnsentry []dnsentry

var dnsdatabase Alldnsentry

func LoadDatabase() []dnsentry {
	file, _ := os.Open("database.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	var configuration Alldnsentry
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func GetIPAdress(dns string, database Alldnsentry) string {
	for _, entry := range database {
		if entry.Dns == dns {
			return entry.Ip
		}
	}
	return "Not exist!"
}
