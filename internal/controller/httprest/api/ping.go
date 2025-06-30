package api

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"
)

type PingAPI struct {
	ping *ping.PingUsecase
}

func NewPingAPI(ping *ping.PingUsecase) *PingAPI {
	return &PingAPI{ping: ping}
}

func (api *PingAPI) Ping(w http.ResponseWriter, r *http.Request) {
	if err := api.ping.Ping(r.Context()); err != nil {
		errInternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
