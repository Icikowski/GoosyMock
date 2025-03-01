package logs_test

import (
	"bytes"
	"testing"

	"git.sr.ht/~icikowski/goosymock/constants"
	"git.sr.ht/~icikowski/goosymock/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmergencyLogger(t *testing.T) {
	shouldContain := []string{constants.ComponentInit, "emergency log"}

	dst := bytes.NewBuffer([]byte{})
	log := logs.GetEmergencyLogger(dst)

	log.Warn().Msg("emergency log")

	output := dst.String()
	hasAll := true
	for _, phrase := range shouldContain {
		hasAll = hasAll && assert.Contains(t, output, phrase)
	}
	require.True(t, hasAll, "some of expected phrases were not found")
}
