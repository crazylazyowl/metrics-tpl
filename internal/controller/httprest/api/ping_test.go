package api

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/api/mocks"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPing_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mocks.NewMockPinger(ctrl)
	usecase := ping.New(repository)
	router := NewPingRouter(usecase)
	server := httptest.NewServer(router)
	client := resty.New().SetBaseURL(server.URL)

	type mock struct {
		err error
	}
	type want struct {
		status int
	}
	tests := []struct {
		name string
		mock mock
		want want
	}{
		{
			"Ping success", mock{err: nil}, want{status: 200},
		},
		{
			"Ping error", mock{err: errors.New("")}, want{status: 500},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.EXPECT().Ping(gomock.Any()).Return(tt.mock.err).Times(1)
			resp, err := client.R().Get("/")
			require.NoError(t, err)
			require.Equal(t, tt.want.status, resp.StatusCode())
		})
	}
}
