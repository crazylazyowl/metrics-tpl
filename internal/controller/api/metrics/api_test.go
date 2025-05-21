package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/stretchr/testify/require"
)

func TestAPI_Update(t *testing.T) {
	repository := memstorage.NewStorage()
	usecase := metrics.NewUsecase(repository)
	api := NewAPI(usecase)

	type args struct {
		method   string
		endpoint string
	}
	type want struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"404 - update counter", args{http.MethodPost, "/update/counter/"}, want{404}},
		{"404 - update gauge", args{http.MethodPost, "/update/gauge/"}, want{404}},
		{"200", args{http.MethodPost, "/update/gauge/testGauge/100"}, want{200}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, tt.args.endpoint, nil)
			rec := httptest.NewRecorder()
			api.Update(rec, req)
			resp := rec.Result()
			defer resp.Body.Close()
			require.Equal(t, tt.want.status, resp.StatusCode)
		})
	}
}
