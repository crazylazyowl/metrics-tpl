package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/api/mocks"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMetrics_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registry := mocks.NewMockMetricRegistry(ctrl)
	usecase := metrics.New(registry)
	router := NewMetricsRouter(usecase)
	server := httptest.NewServer(router)
	client := resty.New().SetBaseURL(server.URL)

	type mock struct {
		err   error
		times int
	}
	type want struct {
		status int
	}
	tests := []struct {
		name     string
		method   string
		endpoint string
		mock     mock
		want     want
	}{
		{
			"Unknown counter metric", http.MethodPost, "/update/counter/",
			mock{nil, 0},
			want{http.StatusNotFound},
		},
		{
			"Unknown gauge metric", http.MethodPost, "/update/gauge/",
			mock{nil, 0},
			want{http.StatusNotFound},
		},
		{
			"Unknown metric type", http.MethodPost, "/update/unknown/testCounter/111",
			mock{nil, 0},
			want{http.StatusBadRequest},
		},
		{
			"Gauge update", http.MethodPost, "/update/gauge/testGauge/100",
			mock{nil, 1},
			want{http.StatusOK},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(tt.mock.err).
				Times(tt.mock.times)

			resp, err := client.R().Execute(tt.method, tt.endpoint)
			require.NoError(t, err)

			require.Equal(t, tt.want.status, resp.StatusCode())
		})
	}
}
