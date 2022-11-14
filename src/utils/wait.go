package utils

import (
	"time"

	"github.com/Icikowski/kubeprobes"
)

// WaitForProbeDown ensures that the probe is down before
// continuing process
func WaitForProbeDown(fn kubeprobes.ProbeFunction) {
	for fn() == nil {
		time.Sleep(500 * time.Millisecond)
	}
}
