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

	c := resty.New()
	c.SetBaseURL(server.URL)

	tests := []struct {
		name   string
		err    error
		status int
	}{
		{
			"Ping success", nil, 200,
		},
		{
			"Ping error", errors.New(""), 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.EXPECT().Ping(gomock.Any()).Return(tt.err).Times(1)
			resp, err := c.R().Get("/")
			require.NoError(t, err)
			require.Equal(t, tt.status, resp.StatusCode())
		})
	}
}
