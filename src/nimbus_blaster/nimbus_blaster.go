/*blaster

  A program to upload large files to nimbus.io in parallel as conjoined archive
*/
package main

import (
	"log"
	"nimbusapi"
	"os"
)

func startConjoined(credentials *nimbusapi.Credentials, flags flags) (
	string, error) {
	requester, err := nimbusapi.NewRequester(credentials)
	if err != nil {
		return "", err
	}

	conjoinedIdentifier, err := nimbusapi.StartConjoined(requester,
		flags.collection, flags.key)
	if err != nil {
		return "", err
	}

	return conjoinedIdentifier, nil
}
func finishConjoined(credentials *nimbusapi.Credentials, flags flags,
	conjoinedIdentifier string) error {
	requester, err := nimbusapi.NewRequester(credentials)
	if err != nil {
		return err
	}

	err = nimbusapi.FinishConjoined(requester, flags.collection, flags.key,
		conjoinedIdentifier)
	if err != nil {
		return err
	}

	return nil
}
func abortConjoined(credentials *nimbusapi.Credentials, flags flags,
	conjoinedIdentifier string) {
	requester, err := nimbusapi.NewRequester(credentials)
	if err != nil {
		return
	}

	nimbusapi.AbortConjoined(requester, flags.collection, flags.key,
		conjoinedIdentifier)
}

func main() {
	log.Println("program starts")
	flags, err := loadflags()
	if err != nil {
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

	if flags.collection == "" {
		flags.collection = nimbusapi.DefaultCollectionName(credentials.Name)
	}

	info, err := os.Stat(flags.filePath)
	if err != nil {
		log.Fatalf("Unable to stat %s %s", flags.filePath, err)
	}
	sliceCount := int(info.Size() / flags.sliceSize)
	if info.Size()%flags.sliceSize != 0 {
		sliceCount += 1
	}
	log.Printf("archiving %s %d bytes %d slices", flags.filePath, info.Size(),
		sliceCount)

	conjoinedIdentifier, err := startConjoined(credentials, flags)
	if err != nil {
		log.Fatalf("StartConjoined %s %s failed %s", flags.collection,
			flags.key, err)
	}
	log.Printf("conjoined_identifier = %s", conjoinedIdentifier)

	work := make(chan workUnit, sliceCount)
	results := make(chan workResult, sliceCount)
	for id := 0; id < flags.connectionCount; id++ {
		requester, err := nimbusapi.NewRequester(credentials)
		if err != nil {
			log.Fatalf("Error creating requester %s\n", err)
		}
		go worker(id, flags.filePath, requester, work, results)
	}

	var offset int64
	var size = flags.sliceSize
	for conjoinedPart := 1; conjoinedPart <= sliceCount; conjoinedPart++ {
		if conjoinedPart == sliceCount {
			size = info.Size() - offset
		}
		workUnit := workUnit{
			flags.collection,
			flags.key,
			conjoinedIdentifier,
			conjoinedPart,
			offset,
			size,
		}
		work <- workUnit
		offset += size
	}

	var completedSize int64
	for completed := 0; completed < sliceCount; completed++ {
		workResult := <-results
		if workResult.err != nil {
			abortConjoined(credentials, flags, conjoinedIdentifier)
			log.Fatalf("Error in worker %d %s %s\n", workResult.workerID,
				workResult.err, workResult.action)
		}
		completedSize += workResult.size
		completedPercent := int(
			float64(completedSize) / float64(info.Size()) * 100.0)
		log.Printf("worker %d completed conjoinedPart %d %d%%",
			workResult.workerID, workResult.conjoinedPart, completedPercent)
	}
	close(work)
	close(results)

	err = finishConjoined(credentials, flags, conjoinedIdentifier)
	if err != nil {
		log.Fatalf("FinishConjoined %s %s failed %s", flags.collection,
			flags.key, err)
	}
	log.Printf("archive complete %s %s conjoined_identifier = %s",
		flags.collection, flags.key, conjoinedIdentifier)

	log.Println("program ends")
}
