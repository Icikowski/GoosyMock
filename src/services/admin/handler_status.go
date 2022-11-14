package admin

import (
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/rs/zerolog"
)

func (s *AdminAPIService) statusHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)

	hostname, _ := os.Hostname()
	ifaces := map[string][]string{}

	netIfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range netIfaces {
			ips := []string{}
			addrs, err := iface.Addrs()
			if err == nil {
				for _, addr := range addrs {
					ips = append(ips, addr.String())
				}
			} else {
				log.Warn().Err(err).Str("interface", iface.Name).Msg("could not list addresses of network interface")
			}
			ifaces[iface.Name] = ips
		}
	} else {
		log.Warn().Err(err).Msg("could not list network interfaces")
	}

	writeResponse(w, r, http.StatusOK, statusResponse{
		System: systemStatusResponsePart{
			OperatingSystem: runtime.GOOS,
			Architecture:    runtime.GOARCH,
		},
		Runtime: runtimeStatusResponsePart{
			GoVersion:  runtime.Version(),
			CPUs:       runtime.NumCPU(),
			Goroutines: runtime.NumGoroutine(),
		},
		Network: networkStatusResponsePart{
			Hostname:   hostname,
			Interfaces: ifaces,
		},
		Stats: statsResponsePart{
			Routes:   s.routes.Count(),
			Payloads: s.payloads.Count(),
		},
	})
}
