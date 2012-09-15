package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"nimbus.io/nimbusapi"
	"strings"
)

const (
	testKey = "test key"
	testBody = "test body"
)

func main() {
	fmt.Println("start")
	var credentials *nimbusapi.Credentials
	var err error

	sp := flag.String("credentials", "", "path to credentials file")
	flag.Parse()
	if *sp == "" {
		credentials, err = nimbusapi.LoadCredentialsFromDefault()
	} else {
		credentials, err = nimbusapi.LoadCredentialsFromPath(*sp)
	}
	if err != nil {
		log.Fatalf("Error loading credentials %s\n", err)
	}

	requester, err := nimbusapi.NewRequester(credentials); if err != nil {
		log.Fatalf("NewRequester failed %s\n", err)
	}

	collectionList, err := nimbusapi.ListCollections(requester, 
		credentials.Name)
	if err != nil {
		log.Fatalf("Request failed %s\n", err)
	}
	fmt.Printf("starting collection list = %v\n", collectionList)

	collectionName := nimbusapi.ReservedCollectionName(credentials.Name, 
		fmt.Sprintf("test-%05d", len(collectionList)))
	collection, err := nimbusapi.CreateCollection(requester, 
		credentials.Name, collectionName)
	if err != nil{
		log.Fatalf("CreateCollection failed %s\n", err)
	}
	fmt.Printf("created collection = %v\n", collection)

	archiveBody := strings.NewReader(testBody)
	versionIdentifier, err := nimbusapi.Archive(requester, collectionName, 
		testKey, archiveBody)
	if err != nil{
		log.Fatalf("Archive failed %s\n", err)
	}
	fmt.Printf("archived key '%s' to version %v\n", testKey, versionIdentifier)

	retrieveBody, err := nimbusapi.Retrieve(requester, collectionName, 
		testKey)
	if err != nil{
		log.Fatalf("Retrieve failed %s\n", err)
	}

	retrieveResult, err := ioutil.ReadAll(retrieveBody)
	retrieveBody.Close()
	if err != nil{
		log.Fatalf("read failed %s\n", err)
	}
	fmt.Printf("retrieved key '%s'; matches testBody = %v\n", testKey, 
		string(retrieveResult) == testBody)

	keySlice, _, err := nimbusapi.ListKeysInCollection(requester, 
		collectionName)
	if err != nil{
		log.Fatalf("ListKeysInCollection failed %s\n", err)
	}
	for _, keyEntry := range keySlice {
		fmt.Printf("deleting key %v\n", keyEntry)
		err = nimbusapi.DeleteKey(requester, collectionName, keyEntry.Name)
		if err != nil {
			log.Fatalf("DeleteKey %v failed %s\n", keyEntry, err)
		}
	}

	success, err := nimbusapi.DeleteCollection(requester, credentials.Name, 
		collectionName)
	if err != nil{
		log.Fatalf("DeleteCollection failed %s\n", err)
	}
	fmt.Printf("deleted collection = %s %v\n", collectionName, success)

	fmt.Println("end")
}
