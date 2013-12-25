/*blaster

  A program to upload large files to nimbus.io in parallel as conjoined archive
*/
package main

import (
	"io"
	"nimbusapi"
	"os"
)

type workUnit struct {
	collection          string
	key                 string
	conjoinedIdentifier string
	conjoinedPart       int
	offset              int64
	size                int64
}

type workResult struct {
	workerID      int
	conjoinedPart int
	size          int64
	err           error
	action        string
}

func worker(id int, filePath string, requester nimbusapi.Requester,
	work <-chan workUnit, results chan<- workResult) {
	result := workResult{}
	result.workerID = id

	file, err := os.Open(filePath)
	if err != nil {
		result.err = err
		result.action = "Open"
		results <- result
		return
	}
	defer file.Close()

	for workUnit := range work {

		_, err = file.Seek(workUnit.offset, 0)
		if err != nil {
			result.err = err
			result.action = "Seek"
			results <- result
			return
		}

		conjoinedParams := nimbusapi.ConjoinedParams{
			workUnit.conjoinedIdentifier, workUnit.conjoinedPart}

		_, err := nimbusapi.Archive(requester, workUnit.collection,
			workUnit.key, &conjoinedParams, workUnit.size,
			io.LimitReader(file, workUnit.size))

		if err != nil {
			result.err = err
			result.action = "Archive"
			results <- result
			return
		}

		result.conjoinedPart = workUnit.conjoinedPart
		result.size = workUnit.size
		results <- result
	}
}
