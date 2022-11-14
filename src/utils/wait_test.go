package utils_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Icikowski/GoosyMock/utils"
	"github.com/stretchr/testify/require"
)

func TestWaitForProbeDown(t *testing.T) {
	status, started := atomic.Bool{}, false
	status.Store(false)

	probeFn := func() error {
		if !started {
			go func() {
				time.Sleep(3 * time.Second)
				status.Store(true)
			}()
			started = false
		}
		if status.Load() {
			return fmt.Errorf("DOWN")
		}
		return nil
	}

	start := time.Now()
	utils.WaitForProbeDown(probeFn)
	end := time.Now()

	// Additional 2 seconds for test operations
	require.WithinDuration(t, end, start, 5*time.Second)
}
