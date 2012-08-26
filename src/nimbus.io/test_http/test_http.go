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
	fmt.Printf("starting collection list = %v\n", collectionList)

	collectionName := nimbusiohttp.ReservedCollectionName(credentials.Name, 
		fmt.Sprintf("test-%05d", len(collectionList)))
	collection, err := nimbusiohttp.CreateCollection(requester, credentials, 
		collectionName)
	if err != nil{
		log.Fatalf("CreateCollection failed %s\n", err)
	}
	fmt.Printf("created collection = %v\n", collection)

	success, err := nimbusiohttp.DeleteCollection(requester, credentials, 
		collectionName)
	if err != nil{
		log.Fatalf("DeleteCollection failed %s\n", err)
	}
	fmt.Printf("deleted collection = %s %v\n", collectionName, success)

	fmt.Println("end")
}
