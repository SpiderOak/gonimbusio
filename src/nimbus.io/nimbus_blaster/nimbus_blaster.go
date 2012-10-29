package main 

import (
	"log"
	"nimbus.io/nimbusapi"
	"os"
)

func main() {
	log.Println("program starts")
	flags, err := loadFlags(); if err != nil {
		log.Fatalf("Unable to load flags: %s\n", err)
	}
	log.Printf("flags = %v", flags)

	var credentials *nimbusapi.Credentials
	if flags.credentialsPath == "" {
		credentials, err = nimbusapi.LoadCredentialsFromDefault()
	} else {
		credentials, err = nimbusapi.LoadCredentialsFromPath(
			flags.credentialsPath)
	}
	if err != nil {
		log.Fatalf("Error loading credentials %s\n", err)
	}

	requester, err := nimbusapi.NewRequester(credentials); if err != nil {
		log.Fatalf("Error creating requester %s\n", err)
	}

	conjoinedIdentifier, err := nimbusapi.StartConjoined(requester, 
		flags.collection, flags.key); if err != nil {
		log.Fatalf("StartConjoined %s %s failed %s", flags.collection, 
			flags.key, err)
	}
	log.Printf("conjoined_identifier = %s", conjoinedIdentifier)

	if flags.collection == "" {
		flags.collection = nimbusapi.DefaultCollectionName(credentials.Name)
	}

	info, err := os.Stat(flags.filePath); if err != nil {
		log.Fatalf("Unable to stat %s %s", flags.filePath, err)
	}
	sliceCount := info.Size() / flags.sliceSize
	if info.Size() % flags.sliceSize != 0 {
		sliceCount += 1
	}
	log.Printf("archiving %s %d bytes %d slices", flags.filePath, info.Size(), 
		sliceCount)

	file, err := os.Open(flags.filePath); if err != nil {
		log.Fatalf("error %s opening %s", err, flags.filePath)
	}
	defer file.Close()

	err = nimbusapi.FinishConjoined(requester, flags.collection, flags.key, 
		conjoinedIdentifier); if err != nil {
		log.Fatalf("FinishConjoined %s %s failed %s", flags.collection, 
			flags.key, err)
	}
	log.Printf("archive complete %s %s conjoined_identifier = %s", 
		flags.collection, flags.key, conjoinedIdentifier)

	log.Println("program ends")	
}