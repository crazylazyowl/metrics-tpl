package main

import (
	"log"
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest"
	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

func main() {
	args, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	storage := memstorage.New()

	usecase := metrics.New(storage)

	router := httprest.NewRouter(usecase)

	_ = http.ListenAndServe(args.address, router)
}
