package utils_test

import (
	"testing"

	"github.com/Icikowski/GoosyMock/utils"
	"github.com/stretchr/testify/require"
)

func TestPointerTo(t *testing.T) {
	tests := map[string]struct {
		value any
	}{
		"string": {value: "test"},
		"int":    {value: 123},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			actual := utils.PointerTo(tc.value)
			require.Equal(t, tc.value, *actual)
		})
	}
}

func TestValueOf(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		stringPtr, intPtr := new(string), new(int)
		*stringPtr, *intPtr = "test", 123

		t.Run("string", func(t *testing.T) {
			actual := utils.ValueOf(stringPtr)
			require.Equal(t, actual, *stringPtr)
		})
		t.Run("int", func(t *testing.T) {
			actual := utils.ValueOf(intPtr)
			require.Equal(t, actual, *intPtr)
		})
	})
	t.Run("without value", func(t *testing.T) {
		stringPtr, intPtr := new(string), new(int)
		stringPtr, intPtr = nil, nil

		t.Run("string", func(t *testing.T) {
			actual := utils.ValueOf(stringPtr)
			require.Equal(t, actual, "")
		})
		t.Run("int", func(t *testing.T) {
			actual := utils.ValueOf(intPtr)
			require.Equal(t, actual, 0)
		})
	})
}

func TestValueOrFallback(t *testing.T) {
	stringFallback, intFallback := "foo", 789

	t.Run("with value", func(t *testing.T) {
		stringPtr, intPtr := new(string), new(int)
		*stringPtr, *intPtr = "test", 123

		t.Run("string", func(t *testing.T) {
			actual := utils.ValueOrFallback(stringPtr, stringFallback)
			require.Equal(t, actual, *stringPtr)
		})
		t.Run("int", func(t *testing.T) {
			actual := utils.ValueOrFallback(intPtr, intFallback)
			require.Equal(t, actual, *intPtr)
		})
	})
	t.Run("without value", func(t *testing.T) {
		stringPtr, intPtr := new(string), new(int)
		stringPtr, intPtr = nil, nil

		t.Run("string", func(t *testing.T) {
			actual := utils.ValueOrFallback(stringPtr, stringFallback)
			require.Equal(t, actual, stringFallback)
		})
		t.Run("int", func(t *testing.T) {
			actual := utils.ValueOrFallback(intPtr, intFallback)
			require.Equal(t, actual, intFallback)
		})
	})
}
