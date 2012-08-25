package main

import (
	"flag"
	"fmt"
	"log"
	nimbusiohttp "nimbus.io/http"
)

func main() {
	fmt.Println("start")
	var credentials *nimbusiohttp.Credentials
	var err error

	sp := flag.String("credentials", "", "path to credentials file")
	flag.Parse()
	if *sp == "" {
		credentials, err = nimbusiohttp.LoadCredentialsFromDefault()
	} else {
		credentials, err = nimbusiohttp.LoadCredentialsFromPath(*sp)
	}
	if err != nil {
		log.Fatalf("Error loading credentials %s\n", err)
	}

	requester, err := nimbusiohttp.NewRequester(credentials); if err != nil {
		log.Fatalf("NewRequester failed %s\n", err)
	}

	collectionList, err := nimbusiohttp.ListCollections(requester, credentials)
	if err != nil {
		log.Fatalf("Request failed %s\n", err)
	}
	fmt.Printf("response = %v\n", collectionList)

	fmt.Println("end")
}
