package model

import (
	"net/http"

	"github.com/Icikowski/GoosyMock/constants"
	"github.com/rs/zerolog"
)

// Route represents the route specification; it consists of the HTTP
// method-specific responses and default response (embedded)
type Route struct {
	*Response `json:",inline" yaml:",inline"`
	GET       *Response `json:"get,omitempty" yaml:"get,omitempty"`
	POST      *Response `json:"post,omitempty" yaml:"post,omitempty"`
	PUT       *Response `json:"put,omitempty" yaml:"put,omitempty"`
	PATCH     *Response `json:"patch,omitempty" yaml:"patch,omitempty"`
	DELETE    *Response `json:"delete,omitempty" yaml:"delete,omitempty"`
}

// GetPayloadsIDs returns a map of payload IDs assigned to method names;
// data is returned only if response for given method has payload ID
// assigned
func (r Route) GetPayloadsIDs() map[string]string {
	data := map[string]string{}

	if r.Response != nil && r.Response.PayloadID != nil {
		data[constants.ResponseDefault] = *r.Response.PayloadID
	}
	if r.GET != nil && r.GET.PayloadID != nil {
		data[http.MethodGet] = *r.GET.PayloadID
	}
	if r.POST != nil && r.POST.PayloadID != nil {
		data[http.MethodPost] = *r.POST.PayloadID
	}
	if r.PUT != nil && r.PUT.PayloadID != nil {
		data[http.MethodPut] = *r.PUT.PayloadID
	}
	if r.PATCH != nil && r.PATCH.PayloadID != nil {
		data[http.MethodPatch] = *r.PATCH.PayloadID
	}
	if r.DELETE != nil && r.DELETE.PayloadID != nil {
		data[http.MethodDelete] = *r.DELETE.PayloadID
	}

	return data
}

// GetResponseForMethod returns the response for given HTTP method
// of default response if method-specific response is not available
// and default response is available, otherwise it returns nil
func (r Route) GetResponseForMethod(method string) *Response {
	responses := map[string]*Response{
		http.MethodGet:    r.GET,
		http.MethodPost:   r.POST,
		http.MethodPut:    r.PUT,
		http.MethodPatch:  r.PATCH,
		http.MethodDelete: r.DELETE,
	}

	if response, ok := responses[method]; ok && response != nil {
		return response
	}
	return r.Response
}

// MarshalZerologObject implements zerolog.LogObjectMarshaler
func (r Route) MarshalZerologObject(e *zerolog.Event) {
	if r.GET != nil {
		e.Object("get", *r.GET)
	}
	if r.POST != nil {
		e.Object("post", *r.POST)
	}
	if r.PUT != nil {
		e.Object("put", *r.PUT)
	}
	if r.PATCH != nil {
		e.Object("patch", *r.PATCH)
	}
	if r.DELETE != nil {
		e.Object("delete", *r.DELETE)
	}
	if r.Response != nil {
		e.Object("default", *r.Response)
	}
}
