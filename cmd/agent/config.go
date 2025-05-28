package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

type config struct {
	address        string
	reportInterval int
	pollInterval   int
}

func loadConfig() (conf *config, err error) {
	conf = &config{}

	flag.StringVar(&conf.address, "a", "localhost:8080", "")
	flag.IntVar(&conf.reportInterval, "r", 10, "")
	flag.IntVar(&conf.pollInterval, "p", 10, "")
	flag.Parse()

	if value := os.Getenv("ADDRESS"); value != "" {
		conf.address = value
	}

	if value := os.Getenv("REPORT_INTERVAL"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		conf.reportInterval = n
	}

	if value := os.Getenv("POLL_INTERVAL"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		conf.pollInterval = n
	}

	switch {
	case conf.address == "":
		err = errors.New("the address is not specified")
	case conf.reportInterval < 1 || conf.reportInterval > 100:
		err = errors.New("the report interval value must be between 1 and 99")
	case conf.pollInterval < 1 || conf.pollInterval > 100:
		err = errors.New("the poll internval value must be between 1 and 99")
	}

	return
}
