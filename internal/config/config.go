package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DBConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret string
}

func LoadConfig() (*Config, error) {
	// Loading .env file
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}

	return config, nil
}
