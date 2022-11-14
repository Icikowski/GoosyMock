package model_test

import (
	"encoding/json"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"
)

type logBuffer struct {
	t         *testing.T
	eventKeys []string
	mux       sync.Mutex
}

func (b *logBuffer) fetchKeys(in map[string]any, prefix []string) []string {
	b.t.Helper()

	keys := []string{}
	for key := range in {
		path := append(prefix, key)
		keys = append(keys, strings.Join(path, "."))

		if sub, ok := in[key].(map[string]any); ok {
			keys = append(keys, b.fetchKeys(sub, path)...)
		}
	}

	return keys
}

func (b *logBuffer) getLast() []string {
	return b.eventKeys
}

func (b *logBuffer) Write(p []byte) (int, error) {
	b.t.Helper()
	b.mux.Lock()
	defer b.mux.Unlock()

	var data map[string]any
	json.Unmarshal(p, &data)

	b.eventKeys = b.fetchKeys(data, []string{})

	return len(p), nil
}

var _ io.Writer = &logBuffer{}

func getTestLog(t *testing.T) (zerolog.Logger, *logBuffer) {
	t.Helper()

	lb := &logBuffer{
		t:         t,
		eventKeys: []string{},
		mux:       sync.Mutex{},
	}

	log := zerolog.New(lb)

	return log, lb
}
