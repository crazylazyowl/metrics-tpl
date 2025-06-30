package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

type config struct {
	address string
	storage struct {
		restore        bool
		backupInterval int
		backupPath     string
	}
}

func (conf *config) Validate() (err error) {
	switch {
	case conf.address == "":
		err = errors.New("the address is not specified")
	case conf.storage.backupInterval == 0:
		err = errors.New("the store interval value can't be 0")
	case conf.storage.backupPath == "":
		err = errors.New("the storage backup filepath is not specified")
	}
	return
}

func loadConfig() (conf *config, err error) {
	conf = &config{}

	flag.StringVar(&conf.address, "a", "localhost:8080", "")
	flag.IntVar(&conf.storage.backupInterval, "i", 300, "")
	flag.StringVar(&conf.storage.backupPath, "f", "dump.json", "")
	flag.BoolVar(&conf.storage.restore, "r", false, "")
	flag.Parse()

	if value := os.Getenv("ADDRESS"); value != "" {
		conf.address = value
	}

	if value := os.Getenv("STORE_INTERVAL"); value != "" {
		conf.storage.backupInterval, err = strconv.Atoi(value)
		if err != nil {
			return
		}
	}

	if value := os.Getenv("FILE_STORAGE_PATH"); value != "" {
		conf.storage.backupPath = value
	}

	if value := os.Getenv("RESTORE"); value != "" {
		conf.storage.restore, err = strconv.ParseBool(value)
		if err != nil {
			return
		}
	}

	err = conf.Validate()

	return
}
