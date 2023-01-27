package admin

import (
	"net/http"
	"strings"

	"github.com/Icikowski/GoosyMock/constants"
)

type errorResponse struct {
	Status  int    `json:"status" yaml:"status"`
	Message string `json:"message" yaml:"message"`
	Details string `json:"details" yaml:"details"`
}

func writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	encoder := r.Context().Value(keyEncoder).(func(any) error)
	format := r.Context().Value(keyFormat).(string)

	response := errorResponse{
		Details: err.Error(),
	}

	switch {
	case strings.Contains(response.Details, "already exists"):
		response.Status, response.Message = http.StatusConflict, "resource already exists"
	case strings.Contains(response.Details, "cannot find payload for"):
		response.Status, response.Message = http.StatusBadRequest, "missing payloads"
	case strings.Contains(response.Details, "does not exist"):
		response.Status, response.Message = http.StatusNotFound, "resource does not exist"
	case strings.Contains(response.Details, "json"):
		response.Status, response.Message = http.StatusUnprocessableEntity, "unable to read object from JSON request body"
	case strings.Contains(response.Details, "yaml"):
		response.Status, response.Message = http.StatusUnprocessableEntity, "unable to read object from YAML request body"
	case strings.Contains(response.Details, "payload"):
		response.Status, response.Message = http.StatusBadRequest, "unable to alternate resource"
	case strings.Contains(response.Details, "EOF"):
		response.Status, response.Message = http.StatusUnprocessableEntity, "no content passed"
	case strings.Contains(response.Details, "unsupported media type"):
		response.Status, response.Message = http.StatusUnsupportedMediaType, "unsupported media type"
	default:
		response.Status, response.Message = http.StatusInternalServerError, "unexpected error occurred"
	}

	w.Header().Add(constants.HeaderContentType, format)
	w.WriteHeader(response.Status)
	_ = encoder(response)
}
