package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv     string `envconfig:"APP_ENV" default:"development"`
	Port       int    `envconfig:"PORT" default:"8080"`
	DbHost     string `envconfig:"DB_HOST" default:"localhost"`
	DbPort     int    `envconfig:"DB_PORT" default:"5432"`
	DbName     string `envconfig:"DB_DATABASE" required:"true"`
	DbUser     string `envconfig:"DB_USER" required:"true"`
	DbPassword string `envconfig:"DB_PASSWORD" required:"true"`
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = envconfig.Process("", c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
