package meta_test

import (
	"testing"

	"git.sr.ht/~icikowski/goosymock/meta"
	"github.com/stretchr/testify/require"
)

func TestGetSentByHeader(t *testing.T) {
	// Let's increase code coverage ¯\_(ツ)_/¯
	require.Equal(t, "GoosyMock/unknown", meta.GetSentByHeader())
}
