package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

type Zerolog struct {
	logger zerolog.Logger
}

func NewLogger() *Zerolog {
	return &Zerolog{
		logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Caller().Logger(),
	}
}

func (l *Zerolog) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *Zerolog) Error(msg string, err error) {
	l.logger.Error().Msg(msg)
}

func (l *Zerolog) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}
