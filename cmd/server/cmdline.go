package main

import (
	"errors"
	"flag"
)

type arguments struct {
	hostport string
}

func parseCmdline() (args *arguments, err error) {
	args = &arguments{}
	flag.StringVar(&args.hostport, "a", "localhost:8080", "")
	flag.Parse()
	switch {
	case args.hostport == "":
		err = errors.New("-a is not specified")
	}
	return
}
