package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	minInterval = 1
	maxInterval = 100
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
	case conf.reportInterval < minInterval || conf.reportInterval > maxInterval:
		err = fmt.Errorf("the report interval must be between %d and %d", minInterval, maxInterval)
	case conf.pollInterval < minInterval || conf.pollInterval > maxInterval:
		err = fmt.Errorf("the poll internval must be between %d and %d", minInterval, maxInterval)
	}

	return
}
