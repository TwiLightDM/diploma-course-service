package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Storage struct {
		Endpoint   string
		AccessKey  string
		SecretKey  string
		BucketName string
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

	cfg.Storage.Endpoint = os.Getenv("MINIO_ENDPOINT")
	cfg.Storage.AccessKey = os.Getenv("MINIO_ROOT_USER")
	cfg.Storage.SecretKey = os.Getenv("minio_root_password")
	cfg.Storage.BucketName = os.Getenv("MINIO_BUCKET")

	cfg.GRPCPort = os.Getenv("GRPC_PORT")

	return cfg
}
