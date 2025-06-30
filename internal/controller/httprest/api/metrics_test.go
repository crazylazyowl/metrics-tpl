package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_UpdateMetric(t *testing.T) {
	repository, _ := memstorage.New(context.TODO(), memstorage.Options{
		Restore:        false,
		BackupPath:     "dump.json",
		BackupInterval: time.Duration(1000) * time.Second,
	})
	usecase := metrics.New(repository)
	router := NewMetricsRouter(usecase)
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
		{http.MethodPost, "/update/unknown/testCounter/111", want{http.StatusBadRequest}},
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
