package model

import "github.com/rs/zerolog"

// Response represents the response specification
type Response struct {
	StatusCode  *int               `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
	Headers     *map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	ContentType *string            `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Content     *string            `json:"content,omitempty" yaml:"content,omitempty"`
	PayloadID   *string            `json:"payloadId,omitempty" yaml:"payloadId,omitempty"`
}

// MarshalZerologObject implements zerolog.LogObjectMarshaler
func (r Response) MarshalZerologObject(e *zerolog.Event) {
	if r.StatusCode != nil {
		e.Int("statusCode", *r.StatusCode)
	}
	if r.Headers != nil {
		h := zerolog.Dict()
		for name, value := range *r.Headers {
			h.Str(name, value)
		}
		e.Dict("headers", h)
	}
	if r.ContentType != nil {
		e.Str("contentType", *r.ContentType)
	}
	if r.Content != nil {
		e.Int("contentLength", len(*r.Content))
	}
	if r.PayloadID != nil {
		e.Str("payloadId", *r.PayloadID)
	}
}
