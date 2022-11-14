package admin

import (
	"net/http"

	"github.com/Icikowski/GoosyMock/model"
	"github.com/rs/zerolog"
)

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

	if err := s.routes.Set(routes); err != nil {
		log.Warn().Err(err).Msg("unable to apply routes")
		writeErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
