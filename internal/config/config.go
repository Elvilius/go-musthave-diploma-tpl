package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	Secret               string
	PollInterval         int
}

func getEnvOrDefaultString(envVar string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(envVar string, defaultValue int) int {
	result := defaultValue
	if value, ok := os.LookupEnv(envVar); ok {
		valueInt, err := strconv.Atoi(value)
		if err == nil {
			result = valueInt
		}
	}
	return result
}

func New() *Config {
	cfg := &Config{
		RunAddress:           getEnvOrDefaultString("RUN_ADDRESS", "localhost:8080"),
		DatabaseURI:          getEnvOrDefaultString("DATABASE_URI", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		AccrualSystemAddress: getEnvOrDefaultString("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:3333"),
		PollInterval:         getEnvOrDefaultInt("POLL_INTERVAL", 3),
	}

	runServerAddress := flag.String("a", cfg.RunAddress, "Server address")
	databaseURI := flag.String("d", cfg.DatabaseURI, "Database URI")
	accrualSystemAddress := flag.String("r", cfg.AccrualSystemAddress, "Accrual System Address")

	flag.Parse()

	cfg.RunAddress = *runServerAddress
	cfg.DatabaseURI = *databaseURI
	cfg.AccrualSystemAddress = *accrualSystemAddress

	cfg.Secret = "Secret"
	fmt.Println("Server Address:", cfg.RunAddress)
	fmt.Println("Database URI", cfg.DatabaseURI)
	fmt.Println("ACCRUAL_SYSTEM_ADDRESS:", cfg.AccrualSystemAddress)

	return cfg
}
