package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost     string
	DbPort     uint64
	DbUser     string
	DbPassword string
	DbName     string

	SecretKey string
}

func New() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return &Config{}, nil
	}

	dbPort, err := strconv.ParseUint(os.Getenv("DB_PORT"), 10, 32)
	if err != nil {
		return &Config{}, nil
	}

	return &Config{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     dbPort,
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		SecretKey:  os.Getenv("SECRET_KEY"),
	}, nil
}
