package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func LoadConfig() *Config {
	environment := os.Getenv("ENVIROMENT")

	if environment == "" || environment == "development" {
		// Se estiver em desenvolvimento, carrega do arquivo .env
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Erro ao carregar o arquivo .env")
		}
	}

	return &Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
	}
}
