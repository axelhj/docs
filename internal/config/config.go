package config

type Config struct {
	CouchHost     string
	CouchPort     string
	CouchUser     string
	CouchPassword string
	DbName        string
	AppPort       string
}

func Load() *Config {
	// Read the .env file/s
	return &Config{}
}
