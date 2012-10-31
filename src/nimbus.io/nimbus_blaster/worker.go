package main 

import (
	"io"
	"nimbus.io/nimbusapi"
	"os"
)

type WorkUnit struct {
	collection string
	key string
	conjoinedIdentifier string
	conjoinedPart int
	offset int64
	size int64
}

func worker(id int, filePath string, requester nimbusapi.Requester, 
	work <-chan WorkUnit, results chan<- error) {

	file, err := os.Open(filePath); if err != nil {
		results <- err
		return
	}
	defer file.Close()

	for workUnit := range work {

		_, err = file.Seek(workUnit.offset, 0); if err != nil {
			results <- err
			return
		}

		conjoinedParams := nimbusapi.ConjoinedParams{
			workUnit.conjoinedIdentifier, workUnit.conjoinedPart}

		_, err := nimbusapi.Archive(requester, workUnit.collection, 
			workUnit.key, &conjoinedParams, io.LimitReader(file, workUnit.size))

		results <- err
	}
}
