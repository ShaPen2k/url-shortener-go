package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" envDefault:"local"`
	StoragePath string `yaml:"storage_path" envRequired:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" envDefault:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" envDefault:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" envDefault:"60s"`
	User        string        `yaml:"user" envRequired:"true"`
	Password    string        `yaml:"password" envRequired:"true" env:"HTTP_SERVER_PASSWORD"`
}

func ConfigLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
