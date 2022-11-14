package logs_test

import (
	"bytes"
	"testing"

	"github.com/Icikowski/GoosyMock/config"
	"github.com/Icikowski/GoosyMock/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerFactory(t *testing.T) {
	const (
		warnLevelMsg  string = "warn level"
		infoLevelMsg  string = "info level"
		debugLevelMsg string = "debug level"
		traceLevelMsg string = "trace level"
	)

	tests := map[string]struct {
		cfg           config.LoggingConfig
		shouldContain []string
	}{
		"unknown log level": {
			cfg: config.LoggingConfig{
				Level: "foo",
			},
			shouldContain: []string{
				"unknown log level, falling back to 'info'",
				warnLevelMsg, infoLevelMsg,
			},
		},
		"pretty log": {
			cfg: config.LoggingConfig{
				Pretty: true,
				Level:  "info",
			},
			shouldContain: []string{
				"WRN", "INF",
				warnLevelMsg, infoLevelMsg,
			},
		},
		"standard log level": {
			cfg: config.LoggingConfig{
				Level: "info",
			},
			shouldContain: []string{
				warnLevelMsg, infoLevelMsg,
			},
		},
		"custom log level": {
			cfg: config.LoggingConfig{
				Level: "trace",
			},
			shouldContain: []string{
				warnLevelMsg, infoLevelMsg,
				debugLevelMsg, traceLevelMsg,
			},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			dst := bytes.NewBuffer([]byte{})

			log := logs.NewLoggerFactory(tc.cfg, dst).InstanceFor("test")
			log.Warn().Msg(warnLevelMsg)
			log.Info().Msg(infoLevelMsg)
			log.Debug().Msg(debugLevelMsg)
			log.Trace().Msg(traceLevelMsg)

			output := dst.String()
			hasAll := true
			for _, phrase := range tc.shouldContain {
				hasAll = hasAll && assert.Contains(t, output, phrase)
			}

			require.True(t, hasAll, "some of expected phrases were not found")
		})
	}
}
