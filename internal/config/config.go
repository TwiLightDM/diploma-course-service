package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres struct {
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

	Mongo struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
	}

	GRPCPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file didn't found")
	}

	cfg := &Config{}

	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Name = os.Getenv("POSTGRES_DB")

	cfg.Storage.Endpoint = os.Getenv("MINIO_ENDPOINT")
	cfg.Storage.AccessKey = os.Getenv("MINIO_ROOT_USER")
	cfg.Storage.SecretKey = os.Getenv("minio_root_password")
	cfg.Storage.BucketName = os.Getenv("MINIO_BUCKET")

	cfg.Mongo.Host = os.Getenv("MONGO_HOST")
	cfg.Mongo.Port = os.Getenv("MONGO_PORT")
	cfg.Mongo.User = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	cfg.Mongo.Name = os.Getenv("MONGO_DB")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	cfg.GRPCPort = os.Getenv("GRPC_PORT")

	return cfg
}
