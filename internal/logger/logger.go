package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(env string) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Development: Pretty console output
	if env == "development" {
		log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Production: JSON output
	return &Logger{log}
}

func (l *Logger) Sync() {
	// Zerolog doesn't need explicit sync
}
