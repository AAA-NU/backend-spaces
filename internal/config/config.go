package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var cfg config

type config struct {
	Server  serverConfig
	Storage storageConfig
}

type serverConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type storageConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Username string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	Name     string `env:"DB_NAME" envDefault:"postgres"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	SslMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

func Config() config {
	return cfg
}

func init() {
	// FIXME: TEMPORARY SOLUTION
	godotenv.Load(".env")
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}
}
