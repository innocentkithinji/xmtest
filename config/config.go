package config

import (
	"log"

	"github.com/spf13/viper"
)

var defaults = map[string]string{
	"port":         "8080",
	"service_name": "xm-company-api",
	"environment":  "QA",

	"SIGNING_SECRET":      "secret",
	"KAFKA_URI":           "localhost:9092",
	"COMPANY_COLLECTION":  "companies",
	"USERS_COLLECTION":    "users",
	"DB_NAME":             "xmtest",
	"MONGODB_URI":         "mongodb://root:toor123@localhost:27017/?retryWrites=true&w=majority",
	"COMPANY_EVENT_TOPIC": "company",
}

func init() {
	log.Println("Initializing Config")
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
}

func InitializeConfig() {
	viper.SetEnvPrefix("XMC")
	viper.AutomaticEnv()
}
