package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	GRPCPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file didn't found")
	}

	cfg := &Config{}

	cfg.DB.Host = os.Getenv("POSTGRES_HOST")
	cfg.DB.Port = os.Getenv("POSTGRES_PORT")
	cfg.DB.User = os.Getenv("POSTGRES_USER")
	cfg.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.DB.Name = os.Getenv("POSTGRES_DB")

	cfg.GRPCPort = os.Getenv("GRPC_PORT")

	return cfg
}
