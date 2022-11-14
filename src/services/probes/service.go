package probes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Icikowski/kubeprobes"
	"github.com/rs/zerolog"
)

// ProbesService represents the service for health probes
type ProbesService struct {
	log           zerolog.Logger
	port          int
	appProbe      *kubeprobes.StatefulProbe
	adminApiProbe *kubeprobes.StatefulProbe
	contentProbe  *kubeprobes.StatefulProbe

	probes http.Handler
}

// NewProbesService creates new Probe service instance
func NewProbesService(
	log zerolog.Logger,
	port int,
	appProbe, adminApiProbe, contentProbe *kubeprobes.StatefulProbe,
) *ProbesService {
	return &ProbesService{
		log:           log,
		port:          port,
		appProbe:      appProbe,
		adminApiProbe: adminApiProbe,
		contentProbe:  contentProbe,
		probes: kubeprobes.New(
			kubeprobes.WithLivenessProbes(
				appProbe.GetProbeFunction(),
			),
			kubeprobes.WithReadinessProbes(
				adminApiProbe.GetProbeFunction(),
				contentProbe.GetProbeFunction(),
			),
		),
	}
}

func (s *ProbesService) prepareServer() *http.Server {
	s.log.Debug().Int("port", s.port).Msg("building server")

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.probes,
	}
}

// Run starts the Probes service with given context
func (s *ProbesService) Run(ctx context.Context) {
	s.log.Info().Msg("starting Probes service")
	go func() {
		s.appProbe.MarkAsUp()
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
		s.appProbe.MarkAsDown()
	}()
}
