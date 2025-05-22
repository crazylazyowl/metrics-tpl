package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_Update(t *testing.T) {
	repository := memstorage.NewStorage()
	usecase := metrics.NewUsecase(repository)
	router := NewRouter(usecase)
	server := httptest.NewServer(router)

	type want struct {
		status int
	}
	tests := []struct {
		method   string
		endpoint string
		want     want
	}{
		{http.MethodPost, "/update/counter/", want{http.StatusNotFound}},
		{http.MethodPost, "/update/gauge/", want{http.StatusNotFound}},
		{http.MethodPost, "/update/gauge/testGauge/100", want{http.StatusOK}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = server.URL + tt.endpoint

			resp, err := req.Send()
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, resp.StatusCode())
		})
	}
}
