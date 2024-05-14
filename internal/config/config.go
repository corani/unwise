package config

import (
	"os"
	"time"

	env "github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	dotenv "github.com/joho/godotenv"
)

type Config struct {
	LogLevel string `env:"LOGLEVEL" envDefault:"info"`
	RestAddr string `env:"REST_ADDR" envDefault:":3123"`
	RestPath string `env:"REST_PATH" envDefault:"/api/v2"`
	Token    string `env:"TOKEN"`
	Logger   *log.Logger
}

func MustLoad() *Config {
	conf, err := Load()
	if err != nil {
		panic(err)
	}

	return conf
}

func Load() (*Config, error) {
	conf := new(Config)

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           log.DebugLevel,
	})

	if err := dotenv.Load(); err != nil {
		logger.Errorf("failed to load .env: %v", err)
	}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	if v, err := log.ParseLevel(conf.LogLevel); err == nil {
		logger.SetLevel(v)
	}

	conf.Logger = logger

	if conf.Token == "" {
		conf.Token = uuid.NewString()

		logger.Info("generated new token", "token", conf.Token)
	}

	return conf, nil
}
