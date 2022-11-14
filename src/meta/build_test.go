package meta_test

import (
	"testing"

	"github.com/Icikowski/GoosyMock/meta"
	"github.com/stretchr/testify/require"
)

func TestGetSentByHeader(t *testing.T) {
	// Let's increase code coverage ¯\_(ツ)_/¯
	require.Equal(t, "GoosyMock/unknown", meta.GetSentByHeader())
}
