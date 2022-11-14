package content

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func (s *ContentService) contentLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := s.log.With().Dict("request", zerolog.Dict().
			Str("remote", r.RemoteAddr).
			Str("path", r.URL.Path).
			Str("method", r.Method).
			Bool("secured", r.TLS != nil),
		).Logger()

		log.Debug().Msg("serving request")
		start := time.Now()

		next.ServeHTTP(w, r.WithContext(context.WithValue(
			r.Context(),
			keyLogger,
			log,
		)))

		end := time.Now()
		log.Info().TimeDiff("duration", end, start).Msg("request served")
	})
}
