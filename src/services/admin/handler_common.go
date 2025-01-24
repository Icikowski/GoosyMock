package admin

import (
	"net/http"

	"git.sr.ht/~icikowski/goosymock/constants"
)

func writeResponse(w http.ResponseWriter, r *http.Request, code int, data any) error {
	encoder := r.Context().Value(keyEncoder).(func(any) error)
	format := r.Context().Value(keyFormat).(string)

	w.Header().Set(constants.HeaderContentType, format)
	w.WriteHeader(code)
	return encoder(data)
}

func writePayload(w http.ResponseWriter, payload []byte) error {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeRaw)
	w.Header().Set("Content-Disposition", "attachment")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(payload)
	return err
}

type systemStatusResponsePart struct {
	OperatingSystem string `json:"os" yaml:"os"`
	Architecture    string `json:"arch" yaml:"arch"`
}

type runtimeStatusResponsePart struct {
	GoVersion  string `json:"goVersion" yaml:"goVersion"`
	CPUs       int    `json:"cpus" yaml:"cpus"`
	Goroutines int    `json:"goroutines" yaml:"goroutines"`
}

type networkStatusResponsePart struct {
	Hostname   string              `json:"hostname" yaml:"hostname"`
	Interfaces map[string][]string `json:"interfaces" yaml:"interfaces"`
}

type statsResponsePart struct {
	Routes   int `json:"routes" yaml:"routes"`
	Payloads int `json:"payloads" yaml:"payloads"`
}

type statusResponse struct {
	System  systemStatusResponsePart  `json:"system" yaml:"system"`
	Runtime runtimeStatusResponsePart `json:"runtime" yaml:"runtime"`
	Network networkStatusResponsePart `json:"network" yaml:"network"`
	Stats   statsResponsePart         `json:"stats" yaml:"stats"`
}

type batchUploadResponse struct {
	Ids    map[string]string `json:"ids" yaml:"ids"`
	Errors map[string]string `json:"errors" yaml:"errors"`
}

type batchDeleteResponse struct {
	Ok      bool   `json:"ok" yaml:"ok"`
	Details string `json:"details,omitempty" yaml:"details,omitempty"`
}
