package main

import (
	"errors"
	"flag"
)

type arguments struct {
	hostport       string
	reportInterval int
	pollInterval   int
}

func parseCmdline() (args *arguments, err error) {
	args = &arguments{}
	flag.StringVar(&args.hostport, "a", "localhost:8080", "")
	flag.IntVar(&args.reportInterval, "r", 10, "")
	flag.IntVar(&args.pollInterval, "p", 10, "")
	flag.Parse()
	switch {
	case args.hostport == "":
		err = errors.New("-a is not specified")
	case args.reportInterval < 1 || args.reportInterval > 100:
		err = errors.New("-p value is allowed between 1 and 99")
	case args.pollInterval < 1 || args.pollInterval > 100:
		err = errors.New("-p value is allowed between 1 and 99")
	}
	return
}
