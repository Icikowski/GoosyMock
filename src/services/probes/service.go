package probes

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"pkg.icikowski.pl/kubeprobes"
)

// ProbesService represents the service for health probes
type ProbesService struct {
	log           zerolog.Logger
	addr          string
	appProbe      kubeprobes.ManualProbe
	adminApiProbe kubeprobes.ManualProbe
	contentProbe  kubeprobes.ManualProbe

	probes http.Handler
}

// NewProbesService creates new Probe service instance
func NewProbesService(
	log zerolog.Logger,
	addr string,
	appProbe, adminApiProbe, contentProbe kubeprobes.ManualProbe,
) *ProbesService {
	probes, _ := kubeprobes.New(
		kubeprobes.WithLivenessProbes(appProbe),
		kubeprobes.WithReadinessProbes(adminApiProbe, contentProbe),
	)
	return &ProbesService{
		log:           log,
		addr:          addr,
		appProbe:      appProbe,
		adminApiProbe: adminApiProbe,
		contentProbe:  contentProbe,
		probes:        probes,
	}
}

func (s *ProbesService) prepareServer() *http.Server {
	s.log.Debug().Str("addr", s.addr).Msg("building server")

	return &http.Server{
		Addr:    s.addr,
		Handler: s.probes,
	}
}

// Run starts the Probes service with given context
func (s *ProbesService) Run(ctx context.Context) {
	s.log.Info().Msg("starting Probes service")
	go func() {
		s.appProbe.Pass()
		running := true
		for running {
			server := s.prepareServer()
			crashChan := make(chan error)

			go func() {
				s.log.Debug().Msg("starting server")
				crashChan <- server.ListenAndServe()
			}()

			select {
			case err := <-crashChan:
				s.log.Warn().Err(err).Msg("Probes service crashed, recovering")
			case <-ctx.Done():
				s.log.Info().Msg("closing Probes service")
				running = false
			}

			_ = server.Close()
		}
		s.appProbe.Fail()
	}()
}
