package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	server := New(EnvTest)

	server.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}
