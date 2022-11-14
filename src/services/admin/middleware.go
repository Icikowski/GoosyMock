package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Icikowski/GoosyMock/constants"
	"github.com/Icikowski/GoosyMock/meta"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

func (*AdminAPIService) adminApiHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.HeaderXSentBy, meta.GetSentByHeader())
		next.ServeHTTP(w, r)
	})
}

func (*AdminAPIService) encoderDecoderResolver(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			encoder func(any) error
			decoder func(any) error
			format  string
		)

		outType := r.Header.Get(constants.HeaderAccept)
		inType := r.Header.Get(constants.HeaderContentType)

		switch {
		case strings.Contains(outType, constants.ContentTypeYAML):
			encoder = yaml.NewEncoder(w).Encode
			format = constants.ContentTypeYAML
		default:
			encoder = json.NewEncoder(w).Encode
			format = constants.ContentTypeJSON
		}

		switch {
		case strings.Contains(inType, constants.ContentTypeYAML):
			decoder = yaml.NewDecoder(r.Body).Decode
		case strings.Contains(inType, constants.ContentTypeJSON):
			decoder = json.NewDecoder(r.Body).Decode
		default:
			decoder = func(_ any) error {
				return fmt.Errorf("unsupported media type: %s", inType)
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, keyEncoder, encoder)
		ctx = context.WithValue(ctx, keyDecoder, decoder)
		ctx = context.WithValue(ctx, keyFormat, format)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *AdminAPIService) contentLogger(next http.Handler) http.Handler {
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
			keyLogger, log,
		)))

		end := time.Now()
		log.Info().TimeDiff("duration", end, start).Msg("request served")
	})
}

func (s *AdminAPIService) payloadContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payloadId := chi.URLParam(r, keyPayloadId)
		defer s.log.Debug().Str(keyPayloadId, payloadId).Msg("resolved payload ID for further processing")

		next.ServeHTTP(w, r.WithContext(context.WithValue(
			r.Context(),
			keyPayloadId, payloadId,
		)))
	})
}
