package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Environment string `envconfig:"environment"`
	AppPort     string `envconfig:"app_port"`

	TLSCert string `envconfig:"tls_cert"`
	TLSKey  string `envconfig:"tls_key"`

	JWTSecret          string `envconfig:"jwt_secret"`
	TokenExpiryMinutes int    `envconfig:"token_expiry_minutes"`

	CouchbaseHost     string `envconfig:"couchbase_host"`
	CouchbaseUsername string `envconfig:"couchbase_username"`
	CouchbasePassword string `envconfig:"couchbase_password"`
}

func Get() (conf Config) {
	if env, _ := os.LookupEnv("ENVIRONMENT"); env != "PRODUCTION" {
		err := godotenv.Load()
		if err != nil {
			log.Println(err)
		}
	}

	err := envconfig.Process("", &conf)
	if err != nil {
		log.Println(err)
	}
	return
}
