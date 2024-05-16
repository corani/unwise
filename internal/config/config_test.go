package config

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/require"
)

func TestConfig_MustLoad(t *testing.T) {
	// NOTE(daniel): these tests read from environment variables,
	// so we can't run them in parallel.

	t.Run("user token", func(t *testing.T) {
		rq := require.New(t)

		t.Setenv("TOKEN", "my-token")

		conf := MustLoad()

		rq.Equal("my-token", conf.Token)
	})

	t.Run("no token", func(t *testing.T) {
		rq := require.New(t)

		conf := MustLoad()

		rq.NotEmpty(conf.Token)
	})

	t.Run("default values", func(t *testing.T) {
		rq := require.New(t)

		conf := MustLoad()

		rq.Equal(":3123", conf.RestAddr)
		rq.Equal("/api/v2", conf.RestPath)
		rq.Equal("info", conf.LogLevel)
		rq.Equal(log.InfoLevel, conf.Logger.GetLevel())
	})
}

func TestConfig_PrintBanner(t *testing.T) {
	rq := require.New(t)
	buf := new(bytes.Buffer)

	conf := MustLoad()
	conf.Logger.SetOutput(buf)

	conf.PrintBanner()

	rq.Contains(buf.String(), "version=")
}
