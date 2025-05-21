package main

import (
	"net/http"

	metricsAPI "github.com/crazylazyowl/metrics-tpl/internal/controller/api/metrics"
	"github.com/crazylazyowl/metrics-tpl/internal/controller/middleware"
	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	metricsUsecase "github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

func main() {
	storage := memstorage.NewStorage()

	usecase := metricsUsecase.NewUsecase(storage)

	api := metricsAPI.NewAPI(usecase)

	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/update/",
		middleware.Methods([]string{http.MethodPost},
			// 		middleware.ContentType("text/plain",
			http.StripPrefix("/update/", http.HandlerFunc(api.Update))))

	http.ListenAndServe("localhost:8080", mux)
}
