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
	key            string
	rateLimit      int
}

func (conf *config) Validate() (err error) {
	switch {
	case conf.address == "":
		err = errors.New("the address is not specified")
	case conf.reportInterval < minInterval || conf.reportInterval > maxInterval:
		err = fmt.Errorf("the report interval must be between %d and %d", minInterval, maxInterval)
	case conf.pollInterval < minInterval || conf.pollInterval > maxInterval:
		err = fmt.Errorf("the poll internval must be between %d and %d", minInterval, maxInterval)
	case conf.rateLimit < 1:
		err = errors.New("the rate limit must be greater than 0")
	}
	return
}

func loadConfig() (conf *config, err error) {
	conf = &config{}

	flag.StringVar(&conf.address, "a", "localhost:8080", "")
	flag.IntVar(&conf.reportInterval, "r", 10, "")
	flag.IntVar(&conf.pollInterval, "p", 10, "")
	flag.StringVar(&conf.key, "k", "", "")
	flag.IntVar(&conf.rateLimit, "l", 1, "")
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

	if value := os.Getenv("KEY"); value != "" {
		conf.key = value
	}

	if value := os.Getenv("RATE_LIMIT"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		conf.rateLimit = n
	}

	err = conf.Validate()

	return
}
