package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config contains configuration for console web server.
type Config struct {
	DBHost      string `required:"false" envconfig:"DB_HOST"`
	DBPort      string `required:"false" envconfig:"DB_PORT"`
	DBUser      string `required:"false" envconfig:"DB_USER"`
	DBName      string `required:"false" envconfig:"DB_NAME"`
	DBPassword  string `required:"false" envconfig:"DB_PASSWORD"`
	DatabaseUrl string `required:"false" envconfig:"DATABASE_URL"`
	ServicePort string `required:"false" envconfig:"SERVICE_PORT"`
	ForexAPIUrl string `required:"true" envconfig:"FOREX_API_URL"`
	ForexAPIKey string `required:"true" envconfig:"FOREX_API_KEY"`
}

// InitConfigs gets environment variables needed to run the app
func InitConfigs() (*Config, error) {
	var cfg Config
	err := godotenv.Load("././.env")
	if err != nil {
		log.Println("Error loading .env file, falling back to cli passed env")
	}
	err = envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading environment variables: %v", err)
	}

	if err != nil {
		return nil, err
	}

	log.Println("InitConfigs: ", cfg)

	return &cfg, nil
}
