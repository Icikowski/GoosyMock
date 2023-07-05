package model_test

import (
	"net/http"
	"testing"

	"github.com/Icikowski/GoosyMock/constants"
	"github.com/Icikowski/GoosyMock/model"
	"github.com/Icikowski/limbo/generics"
	"github.com/stretchr/testify/require"
)

func TestResponseMarshalZerologObject(t *testing.T) {
	tests := map[string]struct {
		response     model.Response
		expectedKeys []string
	}{
		"status code": {
			response: model.Response{
				StatusCode: generics.Ptr(http.StatusOK),
			},
			expectedKeys: []string{"statusCode"},
		},
		"headers": {
			response: model.Response{
				Headers: &map[string]string{
					"X-Test":   "true",
					"Location": "http://example.com",
				},
			},
			expectedKeys: []string{"headers", "headers.X-Test", "headers.Location"},
		},
		"content type": {
			response: model.Response{
				ContentType: generics.Ptr(constants.ContentTypeJSON),
			},
			expectedKeys: []string{"contentType"},
		},
		"content": {
			response: model.Response{
				Content: generics.Ptr("Hello world!"),
			},
			expectedKeys: []string{"contentLength"},
		},
		"payload ID": {
			response: model.Response{
				PayloadID: generics.Ptr("1234567890"),
			},
			expectedKeys: []string{"payloadId"},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			log, lb := getTestLog(t)
			log.Info().EmbedObject(tc.response).Send()
			keys := lb.getLast()

			require.Subset(t, keys, tc.expectedKeys)
		})
	}
}
