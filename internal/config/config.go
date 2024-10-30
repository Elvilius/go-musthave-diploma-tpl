package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	PollInterval = 5
	NumWorkers   = 3
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	Secret               string
	PollInterval         int
	NumWorkers           int
}

func getEnvOrDefaultString(envVar string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	return defaultValue
}

func New() *Config {
	cfg := &Config{
		RunAddress:           getEnvOrDefaultString("RUN_ADDRESS", "localhost:8080"),
		DatabaseURI:          getEnvOrDefaultString("DATABASE_URI", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		AccrualSystemAddress: getEnvOrDefaultString("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:3333"),
	}

	runServerAddress := flag.String("a", cfg.RunAddress, "Server address")
	databaseURI := flag.String("d", cfg.DatabaseURI, "Database URI")
	accrualSystemAddress := flag.String("r", cfg.AccrualSystemAddress, "Accrual System Address")

	flag.Parse()

	cfg.RunAddress = *runServerAddress
	cfg.DatabaseURI = *databaseURI
	cfg.AccrualSystemAddress = *accrualSystemAddress
	cfg.PollInterval = PollInterval
	cfg.NumWorkers = NumWorkers
	cfg.Secret = "Secret"
	fmt.Println("Server Address:", cfg.RunAddress)
	fmt.Println("Database URI", cfg.DatabaseURI)
	fmt.Println("ACCRUAL_SYSTEM_ADDRESS:", cfg.AccrualSystemAddress)

	return cfg
}
