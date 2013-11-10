package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	ConfigPath      string
	Duration        int
	UserIdentityDir string
	MaxUsers        int
}

const (
	defaultDuration = 300
)

func loadFlags() (*Flags, error) {
	cp := flag.String("config", "", "path to test configuration file")
	dp := flag.Int("duration", 0, "Duration of the test in seconds")
	ip := flag.String("user-identity-dir", "",
		"path to a directory containing user identity files")
	mp := flag.Int("max-users", 0, "maximum number of users, 0 == use all")
	flag.Parse()

	if *cp == "" {
		return nil, fmt.Errorf("You must specify a path to the config file")
	}
	if *dp == 0 {
		*dp = defaultDuration
	}
	if *ip == "" {
		return nil, fmt.Errorf(
			"You must specify a path to the user identity directory")
	}

	return &Flags{*cp, *dp, *ip, *mp}, nil
}
