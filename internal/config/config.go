package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" envDefault:"local"`
	Application struct {
		Name string `yaml:"name"`
	} `yaml:"application"`

	Server struct {
		Port        string        `yaml:"port" env-default:"8080"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST"`
		Port     string `yaml:"port" env:"POSTGRES_PORT"`
		Username string `yaml:"user" env:"POSTGRES_USER"`
		Password string `yaml:"pass" env:"POSTGRES_PASSWORD"`
		DbName   string `yaml:"db-name" env:"POSTGRES_DB"`
	} `yaml:"database" env-required:"true"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
