package logging

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// New builds a zerolog logger configured for the provided environment.
func New(env string) zerolog.Logger {
	level := zerolog.InfoLevel
	switch strings.ToLower(env) {
	case "debug", "development", "local":
		level = zerolog.DebugLevel
	case "test":
		level = zerolog.WarnLevel
	case "production", "prod":
		level = zerolog.InfoLevel
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	return logger.Level(level)
}
