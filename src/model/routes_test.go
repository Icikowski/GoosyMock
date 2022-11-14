package model_test

import (
	"testing"

	"github.com/Icikowski/GoosyMock/model"
	"github.com/stretchr/testify/require"
)

func TestGetOrderedPaths(t *testing.T) {
	routes := model.Routes{
		"/test":     {},
		"/test/abc": {},
		"/test/ab*": {},
		"/t*":       {},
		"/*":        {},
		"/a":        {},
		"/":         {},
	}

	expected := []string{
		"/*", "/", "/a", "/t*", "/test", "/test/ab*", "/test/abc",
	}

	require.Equal(t, expected, routes.GetOrderedPaths())
}
