package main

import (
	"errors"
	"flag"
	"os"
)

type config struct {
	address string
}

func (conf *config) Validate() (err error) {
	switch {
	case conf.address == "":
		err = errors.New("the address is not specified")
	}
	return
}

func loadConfig() (conf *config, err error) {
	conf = &config{}

	flag.StringVar(&conf.address, "a", "localhost:8080", "")
	flag.Parse()

	if value := os.Getenv("ADDRESS"); value != "" {
		conf.address = value
	}

	err = conf.Validate()

	return
}
