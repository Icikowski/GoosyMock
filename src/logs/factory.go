package logs

import (
	"io"

	"github.com/Icikowski/GoosyMock/config"
	"github.com/rs/zerolog"
)

// LoggerFactory is a factory of loggers used in application
type LoggerFactory struct {
	logger zerolog.Logger
}

// NewLoggerFactory creates a new logger factory basing on given configuration
func NewLoggerFactory(cfg config.LoggingConfig, target io.Writer) *LoggerFactory {
	var logger zerolog.Logger
	var writer io.Writer = target

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
		defer func() {
			logger.Warn().Str("level", cfg.Level).Msg("unknown log level, falling back to 'info'")
		}()
	}

	if cfg.Pretty {
		writer = zerolog.ConsoleWriter{
			Out:        target,
			NoColor:    false,
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	logger = zerolog.New(writer).Level(level).With().Timestamp().Logger()
	return &LoggerFactory{
		logger: logger,
	}
}

// InstanceFor returns a logger instance for the specified component
func (f *LoggerFactory) InstanceFor(component string) zerolog.Logger {
	return f.logger.With().Str("component", component).Logger()
}
