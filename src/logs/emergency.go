package logs

import (
	"io"

	"git.sr.ht/~icikowski/goosymock/config"
	"git.sr.ht/~icikowski/goosymock/constants"
	"github.com/rs/zerolog"
)

// GetEmergencyLogger returns a logger that will be used in case
// of error while starting application
func GetEmergencyLogger(target io.Writer) zerolog.Logger {
	return NewLoggerFactory(config.LoggingConfig{
		Level: "info",
	}, target).InstanceFor(constants.ComponentInit)
}
