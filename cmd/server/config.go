package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

type config struct {
	address         string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func (conf *config) Validate() (err error) {
	switch {
	case conf.address == "":
		err = errors.New("the address is not specified")
	case conf.storeInterval == 0:
		err = errors.New("the store interval value can't be 0")
	case conf.fileStoragePath == "":
		err = errors.New("the storage backup filepath is not specified")
	}
	return
}

func loadConfig() (conf *config, err error) {
	conf = &config{}

	flag.StringVar(&conf.address, "a", "localhost:8080", "")
	flag.IntVar(&conf.storeInterval, "i", 300, "")
	flag.StringVar(&conf.fileStoragePath, "f", "dump.json", "")
	flag.BoolVar(&conf.restore, "r", false, "")
	flag.Parse()

	if value := os.Getenv("ADDRESS"); value != "" {
		conf.address = value
	}

	if value := os.Getenv("STORE_INTERVAL"); value != "" {
		conf.storeInterval, err = strconv.Atoi(value)
		if err != nil {
			return
		}
	}

	if value := os.Getenv("FILE_STORAGE_PATH"); value != "" {
		conf.fileStoragePath = value
	}

	if value := os.Getenv("RESTORE"); value != "" {
		conf.restore, err = strconv.ParseBool(value)
		if err != nil {
			return
		}
	}

	err = conf.Validate()

	return
}
