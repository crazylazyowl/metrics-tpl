package main

import (
	"net/http"

	metricsAPI "github.com/crazylazyowl/metrics-tpl/internal/controller/api/metrics"
	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	metricsUsecase "github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

func main() {
	storage := memstorage.NewStorage()

	usecase := metricsUsecase.NewUsecase(storage)

	router := metricsAPI.NewRouter(usecase)

	_ = http.ListenAndServe("localhost:8080", router)
}
