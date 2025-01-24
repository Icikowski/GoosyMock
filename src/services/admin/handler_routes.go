package admin

import (
	"net/http"
	"strings"

	"git.sr.ht/~icikowski/goosymock/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type multipleErrors []error

func (me multipleErrors) join() error {
	if len(me) == 0 {
		return nil
	}

	strs := []string{}
	for _, e := range me {
		strs = append(strs, e.Error())
	}
	return errors.New(strings.Join(strs, "; "))
}

func (s *AdminAPIService) listRoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes := s.routes.GetAll()
	writeResponse(w, r, http.StatusOK, routes)
}

func (s *AdminAPIService) replaceRoutesHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)
	decoder := r.Context().Value(keyDecoder).(func(any) error)

	var routes model.Routes = model.Routes{}
	defer r.Body.Close()
	if err := decoder(&routes); err != nil {
		log.Warn().Err(err).Msg("unable to parse routes")
		writeErrorResponse(w, r, err)
		return
	}

	errs := multipleErrors{}
	for path, route := range routes {
		for method, payloadID := range route.GetPayloadsIDs() {
			if _, err := s.payloads.Get(payloadID); err != nil {
				log.Warn().Err(err).Str("method", method).Str("route", path).Msg("unable to find payload")
				errs = append(errs, errors.Wrapf(err, "cannot find payload for %s method of '%s' route", method, path))
			}
		}
	}
	if err := errs.join(); err != nil {
		log.Warn().Msg("configuration cannot be applied due to missing payloads")
		writeErrorResponse(w, r, err)
		return
	}

	if err := s.routes.Set(routes); err != nil {
		log.Warn().Err(err).Msg("unable to apply routes")
		writeErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
