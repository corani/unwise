package config

import (
	"os"
	"path/filepath"
	"time"

	env "github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/corani/unwise/cfg"
	"github.com/google/uuid"
	dotenv "github.com/joho/godotenv"
)

type Config struct {
	LogLevel  string `env:"LOGLEVEL" envDefault:"info"`
	RestAddr  string `env:"REST_ADDR" envDefault:":3123"`
	RestPath  string `env:"REST_PATH" envDefault:"/api/v2"`
	DataPath  string `env:"DATA_PATH" envDefault:"/tmp"`
	Token     string `env:"TOKEN"`
	DropTable string `env:"DROP_TABLE" envDefault:"false"`
	Version   string
	Hash      string
	Logger    *log.Logger
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
	conf.Version = cfg.Version()
	conf.Hash = cfg.Hash()

	if conf.Token == "" {
		conf.Token = uuid.NewString()

		logger.Info("generated new token", "token", conf.Token)
	}

	// TODO(daniel): Should we expand '~' as well?
	if v, err := filepath.Abs(conf.DataPath); err == nil {
		conf.DataPath = v
	} else {
		return nil, err
	}

	return conf, nil
}

func (c *Config) PrintBanner() {
	c.Logger.Info("configuration",
		"version", c.Version,
		"hash", c.Hash,
		"logLevel", c.LogLevel,
		"restAddr", c.RestAddr,
		"restPath", c.RestPath,
		"dataPath", c.DataPath,
	)
}
