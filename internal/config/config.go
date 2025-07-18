package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DatabaseURL         string
	RedisURL           string
	EthereumRPC        string
	ContractAddress    string
	PrivateKey         string
	JWTSecret          string
	RateLimitWindow    string
	RateLimitMax       string
}

func Load() (*Config, error) {
	// Carregar .env se existir
	godotenv.Load()

	return &Config{
		Port:            getEnv("API_PORT", "3000"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:password@localhost/vfinance?sslmode=disable"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		EthereumRPC:     getEnv("ETHEREUM_RPC", "http://localhost:8545"),
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		PrivateKey:      getEnv("PRIVATE_KEY", ""),
		JWTSecret:       getEnv("JWT_SECRET", "your_secret_key"),
		RateLimitWindow: getEnv("RATE_LIMIT_WINDOW", "900000"),
		RateLimitMax:    getEnv("RATE_LIMIT_MAX", "100"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}