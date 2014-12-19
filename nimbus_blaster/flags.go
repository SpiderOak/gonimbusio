/*blaster

  A program to upload large files to nimbus.io in parallel as conjoined archive
*/
package main

import (
	"flag"
	"fmt"
)

type flags struct {
	credentialsPath string
	filePath        string
	collection      string
	key             string
	sliceSize       int64
	connectionCount int
}

const (
	defaultSliceSize   = 10 * 1024 * 1024
	defaultConnections = 5
)

func loadflags() (flags, error) {
	credp := flag.String("credentials", "", "path to credentials file")
	fp := flag.String("file-path", "", "path of file to be archived")
	colp := flag.String("collection", "", "collection name")
	keyp := flag.String("key", "", "key name")
	sp := flag.Int64("slices-size", defaultSliceSize,
		"maximum size of each slice of the file to be archived")
	hp := flag.Int("connections", defaultConnections,
		"Max number of open HTTP connections")
	flag.Parse()

	if *fp == "" {
		return flags{}, fmt.Errorf(
			"you must specify the path to a file to be archived")
	}
	if *keyp == "" {
		return flags{}, fmt.Errorf("you must specify a key")
	}
	if *sp == 0 {
		*sp = defaultSliceSize
	}
	if *hp == 0 {
		*hp = defaultConnections
	}

	return flags{*credp, *fp, *colp, *keyp, *sp, *hp}, nil
}
