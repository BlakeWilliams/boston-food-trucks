package static

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blakewilliams/medium/pkg/router"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	testCases := map[string]struct {
		requestPath            string
		expectedResponseBody   string
		expectedResponseStatus int
	}{
		"valid path": {
			requestPath:            "/static/hello.js",
			expectedResponseBody:   "alert(\"the truth is out there\");\n",
			expectedResponseStatus: 200,
		},
		"missing asset": {
			requestPath:            "/static/missing.js",
			expectedResponseBody:   "404 page not found\n",
			expectedResponseStatus: 404,
		},
		"path traversal": {
			requestPath:            "/static/../static.go",
			expectedResponseBody:   `404 not found`,
			expectedResponseStatus: 404,
		},
		"app pass through": {
			requestPath:            "/",
			expectedResponseBody:   `hello from app`,
			expectedResponseStatus: 200,
		},
	}
	r := router.New(router.DefaultActionFactory)

	config := Config{
		FileRoot:   "./test_files",
		PathPrefix: "/static",
	}
	r.Use(Middleware(config))

	r.Get("/", func(ac *router.BaseAction) {
		ac.Response().Write([]byte("hello from app"))
	})

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.requestPath, nil)
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			require.Equal(t, tc.expectedResponseStatus, res.Code)
			require.Equal(t, tc.expectedResponseBody, res.Body.String())
		})
	}
}
