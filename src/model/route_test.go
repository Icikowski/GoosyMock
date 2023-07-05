package model_test

import (
	"net/http"
	"testing"

	"github.com/Icikowski/GoosyMock/constants"
	"github.com/Icikowski/GoosyMock/model"
	"github.com/Icikowski/limbo/generics"
	"github.com/stretchr/testify/require"
)

func TestGetResponseForMethod(t *testing.T) {
	getResponse := &model.Response{
		StatusCode: generics.Ptr(http.StatusOK),
	}
	postResponse := &model.Response{
		ContentType: generics.Ptr(constants.ContentTypeRaw),
	}
	defaultResponse := &model.Response{
		PayloadID: generics.Ptr("some-uuid"),
	}
	routeWithoutDefaultResponse := model.Route{
		GET:  getResponse,
		POST: postResponse,
	}
	routeWithDefaultResponse := model.Route{
		GET:      getResponse,
		Response: defaultResponse,
	}

	tests := map[string]map[string]struct {
		route            model.Route
		expectedResponse *model.Response
	}{
		"default response undefined": {
			http.MethodGet: {
				route:            routeWithoutDefaultResponse,
				expectedResponse: getResponse,
			},
			http.MethodPost: {
				route:            routeWithoutDefaultResponse,
				expectedResponse: postResponse,
			},
			http.MethodPut: {
				route:            routeWithoutDefaultResponse,
				expectedResponse: nil,
			},
		},
		"default response defined": {
			http.MethodGet: {
				route:            routeWithDefaultResponse,
				expectedResponse: getResponse,
			},
			http.MethodPost: {
				route:            routeWithDefaultResponse,
				expectedResponse: defaultResponse,
			},
			http.MethodPut: {
				route:            routeWithDefaultResponse,
				expectedResponse: defaultResponse,
			},
		},
	}

	for name, tcs := range tests {
		name, tcs := name, tcs
		t.Run(name, func(t *testing.T) {
			for method, tc := range tcs {
				method, tc := method, tc
				t.Run(method, func(t *testing.T) {
					actual := tc.route.GetResponseForMethod(method)
					require.Equal(t, tc.expectedResponse, actual)
				})
			}
		})
	}
}

func TestRouteMarshalZerologObject(t *testing.T) {
	someResponse := &model.Response{
		StatusCode: generics.Ptr(http.StatusOK),
	}

	tests := map[string]struct {
		route        model.Route
		expectedKeys []string
	}{
		"get": {
			route: model.Route{
				GET: someResponse,
			},
			expectedKeys: []string{"get.statusCode"},
		},
		"post": {
			route: model.Route{
				POST: someResponse,
			},
			expectedKeys: []string{"post.statusCode"},
		},
		"put": {
			route: model.Route{
				PUT: someResponse,
			},
			expectedKeys: []string{"put.statusCode"},
		},
		"patch": {
			route: model.Route{
				PATCH: someResponse,
			},
			expectedKeys: []string{"patch.statusCode"},
		},
		"delete": {
			route: model.Route{
				DELETE: someResponse,
			},
			expectedKeys: []string{"delete.statusCode"},
		},
		"default": {
			route: model.Route{
				Response: someResponse,
			},
			expectedKeys: []string{"default.statusCode"},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			log, lb := getTestLog(t)
			log.Info().EmbedObject(tc.route).Send()
			keys := lb.getLast()

			require.Subset(t, keys, tc.expectedKeys)
		})
	}
}
