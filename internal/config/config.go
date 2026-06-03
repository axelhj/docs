package config

import (
	"log"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var envFiles = []string{
	".env",
	".env.local",
}

type Config struct {
	CouchHost     string `env:"COUCHDB_HOST"`
	CouchPort     string `env:"COUCHDB_PORT"`
	CouchUser     string `env:"COUCHDB_USER_NAME"`
	CouchPassword string `env:"COUCHDB_USER_PW"`
	DbName        string `env:"DB_NAME"`
	AppPort       string `env:"APP_PORT"`
}

func Load() *Config {
	var err error
	for i, filename := range envFiles {
		err = godotenv.Load(filename)
		if err != nil && i == 0 {
			log.Fatalf("Couldn't load %v: %v\n", filename, err)
		} else if err != nil {
			log.Printf("Couldn't load %v: %v\n", filename, err)
		}
	}
	var config Config
	if err = env.Parse(&config); err != nil {
		log.Fatal("Couldn't parse env into config struct, %v\n", err)
	}
	// Print all of the secrets to stdout
	// log.Printf("%+v\n", config)
	return &config
}
