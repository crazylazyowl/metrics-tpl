package api

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"
)

func NewPingRouter(ping *ping.PingUsecase) http.Handler {
	api := NewPingAPI(ping)

	r := http.NewServeMux()
	r.HandleFunc("/", api.Ping)

	return r
}
