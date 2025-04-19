package serviceauth

import (
	"log"
	"os"
)

type Config struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	AuthServiceURL string
	BootstrapToken string
	ServiceName    string
}

func LoadConfigWithPrefix(prefix string) Config {
	cfg := Config{
		RedisAddr:      os.Getenv(prefix + "REDIS_ADDR"),
		RedisPassword:  os.Getenv(prefix + "REDIS_PASS"),
		AuthServiceURL: os.Getenv(prefix + "AUTH_SERVICE_URL"),
		BootstrapToken: os.Getenv(prefix + "BOOTSTRAP_TOKEN"),
		ServiceName:    os.Getenv(prefix + "SERVICE_NAME"),
	}

	if cfg.RedisAddr == "" || cfg.ServiceName == "" {
		log.Panicf("Missing required environment variables: %sREDIS_ADDR or %sSERVICE_NAME", prefix, prefix)
	}

	return cfg
}
