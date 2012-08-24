package http

import (
	"fmt"
)

type Collection struct {
	name string 
	versioned bool
}

func ListCollections(requester Requester, credentials *Credentials) (
	[]Collection, error) {
	var response *Response
	var err error

	method := "GET"
	baseURI := fmt.Sprintf("/customers/%s/collections", credentials.Name)

	response, err = requester.Request(method, baseURI)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("GET %S failed (%d) %s", baseURI, response.StatusCode, 
			response.Body)
		return nil, err
	}

	return nil, nil

}