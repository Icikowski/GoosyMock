package logs

import (
	"io"

	"github.com/Icikowski/GoosyMock/config"
	"github.com/Icikowski/GoosyMock/constants"
	"github.com/rs/zerolog"
)

// GetEmergencyLogger returns a logger that will be used in case
// of error while starting application
func GetEmergencyLogger(target io.Writer) zerolog.Logger {
	return NewLoggerFactory(config.LoggingConfig{
		Level: "info",
	}, target).InstanceFor(constants.ComponentInit)
}
