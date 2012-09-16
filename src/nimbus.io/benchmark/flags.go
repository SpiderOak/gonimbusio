package main

import (
	"flag"
	"fmt"
)

type flags struct {
	configPath string
	duration int
}

const (
	defaultDuration = 300
)

func loadFlags () (*flags, error) {
	cp := flag.String("config", "", "path to test configuration file")
	flag.Parse()
	dp := flag.Int("duration", 0, "Duration of the test in seconds")
	flag.Parse()

	if *cp == "" {
		return nil, fmt.Errorf("You must specify a path to the config file") 
	}
	if *dp == 0 {
		*dp = defaultDuration
	}

	return &flags{*cp, *dp}, nil
}