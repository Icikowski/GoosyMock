package admin

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (s *AdminAPIService) listPayloadsHandler(w http.ResponseWriter, r *http.Request) {
	payloads := s.payloads.GetAll()
	writeResponse(w, r, http.StatusOK, payloads)
}

func (s *AdminAPIService) uploadPayloadsHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)

	ids, errs := map[string]string{}, map[string]string{}

	if err := r.ParseMultipartForm(s.maxPayloadSize<<20 + 512); err != nil {
		log.Warn().Err(err).Send()
		writeErrorResponse(w, r, errors.Wrap(err, "unable to parse form with payload attachments"))
		return
	}

	payloads, ok := r.MultipartForm.File[formFieldPayloads]
	if !ok {
		err := fmt.Errorf("request does not contain any file in payloads field as expected")
		log.Warn().Err(err).Send()
		writeErrorResponse(w, r, err)
		return
	}
	for _, item := range payloads {
		file, err := item.Open()
		if err != nil {
			log.Warn().Err(err).Msg("unable to process file")
			errs[item.Filename] = err.Error()
			continue
		}

		id, err := s.payloads.InsertFile(item.Filename, file)
		if err != nil {
			log.Warn().Err(err).Send()
			errs[item.Filename] = err.Error()
			continue
		}
		ids[item.Filename] = id
	}

	writeResponse(w, r, http.StatusOK, batchUploadResponse{
		Ids:    ids,
		Errors: errs,
	})
}

func (s *AdminAPIService) deletePayloadsHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)

	if err := s.payloads.DeleteAll(); err != nil {
		log.Warn().Err(err).Send()
		writeResponse(w, r, http.StatusOK, batchDeleteResponse{Details: err.Error()})
		return
	}
	writeResponse(w, r, http.StatusOK, batchDeleteResponse{Ok: true})
}

func (s *AdminAPIService) fetchSinglePayloadHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)
	payloadId := r.Context().Value(keyPayloadId).(string)

	payload, err := s.payloads.Get(payloadId)
	if err != nil {
		writeErrorResponse(w, r, err)
		log.Warn().Err(err).Send()
		return
	}
	writeResponse(w, r, http.StatusOK, payload)
}

func (s *AdminAPIService) downloadSinglePayloadHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)
	payloadId := r.Context().Value(keyPayloadId).(string)

	payload, err := s.payloads.Get(payloadId)
	if err != nil {
		writeErrorResponse(w, r, err)
		log.Warn().Err(err).Send()
		return
	}

	data, err := payload.Contents()
	if err != nil {
		writeErrorResponse(w, r, errors.Wrap(err, "unable to resolve payload contents"))
		log.Warn().Err(err).Send()
		return
	}

	if err := writePayload(w, data); err != nil {
		log.Warn().Err(err).Send()
	}
}

func (s *AdminAPIService) updateSinglePayloadHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)
	payloadId := r.Context().Value(keyPayloadId).(string)

	if err := r.ParseMultipartForm(s.maxPayloadSize<<20 + 512); err != nil {
		log.Warn().Err(err).Send()
		writeErrorResponse(w, r, errors.Wrap(err, "unable to process payload request"))
		return
	}

	items, ok := r.MultipartForm.File["payload"]
	if !ok {
		err := fmt.Errorf("request does not contain any file in payload field as expected")
		log.Warn().Err(err).Send()
		writeErrorResponse(w, r, err)
		return
	}
	if len(items) != 1 {
		err := fmt.Errorf("request does not contain single file in payload field as expected")
		log.Warn().Err(err).Send()
		writeErrorResponse(w, r, err)
		return
	}

	item := items[0]
	file, err := item.Open()
	if err != nil {
		log.Warn().Err(err).Msg("unable to process payload")
		writeErrorResponse(w, r, errors.Wrap(err, "error while processing payload"))
		return
	}

	if err := s.payloads.UpdateFile(payloadId, item.Filename, file); err != nil {
		log.Warn().Err(err).Msg("unable to update payload")
		writeErrorResponse(w, r, errors.Wrap(err, "error while updating payload"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s *AdminAPIService) deleteSinglePayloadHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(keyLogger).(zerolog.Logger)
	payloadId := r.Context().Value(keyPayloadId).(string)

	if err := s.payloads.Delete(payloadId); err != nil {
		log.Warn().Err(err).Msg("unable to delete payload")
		writeErrorResponse(w, r, errors.Wrap(err, "error while deleting payload"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
