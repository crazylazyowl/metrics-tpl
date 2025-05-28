package main

import (
	"errors"
	"flag"
	"os"
)

type config struct {
	address string
}

func loadConfig() (conf *config, err error) {
	conf = &config{}

	flag.StringVar(&conf.address, "a", "localhost:8080", "")
	flag.Parse()

	if value := os.Getenv("ADDRESS"); value != "" {
		conf.address = value
	}

	switch {
	case conf.address == "":
		err = errors.New("the address is not specified")
	}

	return
}
