package main

import (
	"log"
	"net/http"

	metricsAPI "github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/metrics"
	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	metricsUsecase "github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

func main() {
	args, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	storage := memstorage.New()

	usecase := metricsUsecase.New(storage)

	router := metricsAPI.NewRouter(usecase)

	_ = http.ListenAndServe(args.address, router)
}
