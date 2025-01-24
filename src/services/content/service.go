package content

import (
	"context"
	"net/http"

	"git.sr.ht/~icikowski/goosymock/config"
	"git.sr.ht/~icikowski/goosymock/constants"
	"git.sr.ht/~icikowski/goosymock/data"
	"git.sr.ht/~icikowski/goosymock/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"pkg.icikowski.pl/kubeprobes"
)

// ContentService represents the service for content serving
type ContentService struct {
	log   zerolog.Logger
	cfg   config.ServiceConfig
	probe kubeprobes.ManualProbe

	routes   data.SubscribableStore[model.Route]
	payloads data.Store[data.Payload]

	handler http.Handler
}

// NewContentService creates new Content service instance
func NewContentService(
	log zerolog.Logger,
	cfg config.ServiceConfig,
	probe kubeprobes.ManualProbe,
) *ContentService {
	service := &ContentService{
		log:   log,
		cfg:   cfg,
		probe: probe,
	}

	return service
}

const (
	rootPath  string = "/*"
	keyLogger string = "logger"
)

func (s *ContentService) buildHandler(routes model.Routes) {
	s.log.Debug().Msg("building handler")

	handler := chi.NewRouter()
	handler.Use(middleware.RealIP)
	handler.Use(middleware.CleanPath)
	handler.Use(s.contentLogger)
	handler.Use(middleware.Recoverer)

	if _, ok := routes[rootPath]; !ok {
		s.log.Debug().Msg("no default route found, using echo handler")
		handler.HandleFunc(rootPath, defaultHandler)
	}

	for _, path := range routes.GetOrderedPaths() {
		handler.HandleFunc(path, s.getRouteHandler(routes[path]))
	}

	s.handler = handler
}

func (s *ContentService) buildServers() (*http.Server, *http.Server) {
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

// Run starts the Content service with given context; it returns
// a function that can be used for updating the routes configuration
func (s *ContentService) Run(
	ctx context.Context,
	routesStore data.SubscribableStore[model.Route],
	payloadsStore data.Store[data.Payload],
) {
	s.routes = routesStore
	s.payloads = payloadsStore

	s.log.Debug().Msg("subscribing to routes store changes")
	routesChanged := s.routes.Subscribe(constants.ComponentContentService)

	s.log.Info().Msg("starting Content service")
	go func() {
		running := true
		for running {
			s.buildHandler(s.routes.GetAll())
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
			case <-routesChanged:
				s.log.Info().Msg("reloading routes configuration")
				s.buildHandler(s.routes.GetAll())
			case err := <-crashChan:
				s.log.Warn().Err(err).Msg("Content service crashed, recovering")
			case <-ctx.Done():
				s.log.Info().Msg("closing Content service")
				running = false
			}

			_ = plain.Close()
			if secured != nil {
				_ = secured.Close()
			}

			s.probe.Fail()
		}
		s.routes.Unsubscribe(constants.ComponentContentService)
	}()
}
