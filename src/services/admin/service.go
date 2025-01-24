package admin

import (
	"context"
	"net/http"

	"git.sr.ht/~icikowski/goosymock/config"
	"git.sr.ht/~icikowski/goosymock/data"
	"git.sr.ht/~icikowski/goosymock/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"pkg.icikowski.pl/kubeprobes"
)

// AdminAPIService represents the service for application management
type AdminAPIService struct {
	log   zerolog.Logger
	cfg   config.ServiceConfig
	probe kubeprobes.ManualProbe

	maxPayloadSize int64
	routes         data.Store[model.Route]
	payloads       data.FilesStore[data.Payload]

	handler http.Handler
}

// NewAdminAPIService creates new Admin API service
func NewAdminAPIService(
	log zerolog.Logger,
	cfg config.ServiceConfig,
	probe kubeprobes.ManualProbe,
	maxPayloadSize int64,
) *AdminAPIService {
	srv := &AdminAPIService{
		log:            log,
		cfg:            cfg,
		probe:          probe,
		maxPayloadSize: maxPayloadSize,
	}

	srv.buildHandler()

	return srv
}

const (
	rootPath          string = "/"
	formFieldPayloads string = "payloads"
	keyPayloadId      string = "payloadId"
	keyLogger         string = "logger"
	keyEncoder        string = "encoder"
	keyDecoder        string = "decoder"
	keyFormat         string = "format"
)

func (s *AdminAPIService) buildHandler() {
	handler := chi.NewRouter()
	handler.Use(middleware.RealIP)
	handler.Use(middleware.CleanPath)
	handler.Use(s.contentLogger)
	handler.Use(s.adminApiHeader)
	handler.Use(s.encoderDecoderResolver)
	handler.Use(middleware.Recoverer)

	handler.Get(rootPath, s.statusHandler)

	handler.Route("/routes", func(r chi.Router) {
		r.Get(rootPath, s.listRoutesHandler)
		r.Post(rootPath, s.replaceRoutesHandler)
	})

	handler.Route("/payloads", func(r chi.Router) {
		r.Get(rootPath, s.listPayloadsHandler)
		r.Post(rootPath, s.uploadPayloadsHandler)
		r.Delete(rootPath, s.deletePayloadsHandler)

		r.Route("/{payloadId}", func(r chi.Router) {
			r.Use(s.payloadContext)
			r.Get(rootPath, s.fetchSinglePayloadHandler)
			r.Get("/download", s.downloadSinglePayloadHandler)
			r.Post(rootPath, s.updateSinglePayloadHandler)
			r.Delete(rootPath, s.deleteSinglePayloadHandler)
		})
	})

	s.handler = handler
}

func (s *AdminAPIService) buildServers() (*http.Server, *http.Server) {
	s.log.Debug().Dict("addrs", zerolog.Dict().
		Str("plain", s.cfg.Address).
		Str("secured", s.cfg.SecuredAddress),
	).Bool("sslEnabled", s.cfg.TLSEnabled).Msg("building servers")

	plainServer := &http.Server{
		Addr:    s.cfg.Address,
		Handler: s.handler,
	}

	var securedServer *http.Server
	if s.cfg.TLSEnabled {
		securedServer = &http.Server{
			Addr:      s.cfg.SecuredAddress,
			TLSConfig: s.cfg.GetTLSConfig(),
			Handler:   s.handler,
		}
	}

	return plainServer, securedServer
}

// Run starts the Admin API service with given context
func (s *AdminAPIService) Run(
	ctx context.Context,
	routesStore data.Store[model.Route],
	payloadsStore data.FilesStore[data.Payload],
) {
	s.routes = routesStore
	s.payloads = payloadsStore

	s.log.Info().Msg("starting Admin API service")
	go func() {
		running := true
		for running {
			plain, secured := s.buildServers()
			crashChan := make(chan error)

			go func() {
				s.log.Debug().Msg("starting plain server")
				crashChan <- plain.ListenAndServe()
			}()
			if secured != nil {
				s.log.Debug().Msg("starting secured server")
				go func() {
					crashChan <- secured.ListenAndServeTLS("", "")
				}()
			}
			s.probe.Pass()

			select {
			case err := <-crashChan:
				s.log.Warn().Err(err).Msg("Admin API service crashed, recovering")
			case <-ctx.Done():
				s.log.Info().Msg("closing Admin API service")
				running = false
			}

			_ = plain.Close()
			if secured != nil {
				_ = secured.Close()
			}

			s.probe.Fail()
		}
	}()
}
