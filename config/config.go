package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	CorsAllowedOrigins []string
	PostgresHost       string
	PostgresPort       string
	PostgresUser       string
	PostgresPassword   string
	PostgresDB         string

	// Google Oauth provider
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// JWT
	JWTSecret string
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
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		CorsAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "eventdb"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),

		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

func ConnectDB() (*pgxpool.Pool, error) {
	cfg := LoadConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB,
	)

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}
