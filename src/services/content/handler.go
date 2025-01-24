package content

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sr.ht/~icikowski/goosymock/constants"
	"git.sr.ht/~icikowski/goosymock/meta"
	"git.sr.ht/~icikowski/goosymock/model"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"pkg.icikowski.pl/generics"
)

type errorResponse struct {
	Status  int    `json:"status" yaml:"status"`
	Message string `json:"message" yaml:"message"`
	Details string `json:"details" yaml:"details"`
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	response := errorResponse{
		Status:  http.StatusInternalServerError,
		Message: "unable to fetch associated payload",
		Details: err.Error(),
	}

	w.Header().Add(constants.HeaderContentType, constants.ContentTypeJSON)
	w.WriteHeader(response.Status)
	_ = json.NewEncoder(w).Encode(response)
}

func (s *ContentService) getRouteHandler(route model.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value(keyLogger).(zerolog.Logger)

		response := route.GetResponseForMethod(r.Method)
		if response == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		data := []byte(generics.Val(response.Content))
		if response.PayloadID != nil {
			payload, err := s.payloads.Get(*response.PayloadID)
			if err != nil {
				log.Err(err).Str("payloadId", *response.PayloadID).Msg("unable to find payload with given ID")
				writeErrorResponse(w, err)
				return
			}

			data, err = payload.Contents()
			if err != nil {
				log.Err(err).Str("payloadId", *response.PayloadID).Msg("unable to fetch payload with given ID")
				writeErrorResponse(w, err)
				return
			}
		}

		w.Header().Set(
			constants.HeaderContentType,
			generics.Fallback(
				response.ContentType,
				constants.DefaultResponseContentType,
			),
		)

		for key, value := range generics.Val(response.Headers) {
			w.Header().Set(key, value)
		}

		w.WriteHeader(generics.Fallback(
			response.StatusCode,
			constants.DefaultResponseStatusCode,
		))

		w.Write(data)
	}
}

type defaultResponse struct {
	Remote  string              `json:"remote" yaml:"remote"`
	Host    string              `json:"host" yaml:"host"`
	Path    string              `json:"path" yaml:"path"`
	Method  string              `json:"method" yaml:"method"`
	Headers map[string]string   `json:"headers,omitempty" yaml:"headers,omitempty"`
	Queries map[string][]string `json:"queries,omitempty" yaml:"queries,omitempty"`
	Cookies map[string]string   `json:"cookies,omitempty" yaml:"cookies,omitempty"`
	JSON    any                 `json:"json,omitempty" yaml:"json,omitempty"`
	YAML    any                 `json:"yaml,omitempty" yaml:"yaml,omitempty"`
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	headers := map[string]string{}
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ", ")
	}

	queries := map[string][]string{}
	if err := r.ParseForm(); err == nil {
		queries = r.Form
	}

	cookies := map[string]string{}
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	var jsonData, yamlData any
	switch r.Header.Get(constants.HeaderContentType) {
	case constants.ContentTypeJSON:
		_ = json.NewDecoder(r.Body).Decode(&jsonData)
	case constants.ContentTypeYAML:
		_ = yaml.NewDecoder(r.Body).Decode(&yamlData)
	}

	data := defaultResponse{
		Remote:  r.RemoteAddr,
		Host:    r.Host,
		Path:    r.URL.Path,
		Method:  r.Method,
		Headers: headers,
		Queries: queries,
		Cookies: cookies,
		JSON:    jsonData,
		YAML:    yamlData,
	}

	contentType := constants.ContentTypeJSON
	var encoder interface {
		Encode(any) error
	} = json.NewEncoder(w)

	if r.Header.Get(constants.HeaderAccept) == constants.ContentTypeYAML {
		contentType = constants.ContentTypeYAML
		encoder = yaml.NewEncoder(w)
	}

	w.Header().Add(constants.HeaderContentType, contentType)
	w.Header().Add(constants.HeaderXSentBy, meta.GetSentByHeader())
	w.WriteHeader(http.StatusOK)

	encoder.Encode(data)
}
