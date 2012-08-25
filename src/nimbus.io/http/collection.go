package http

import (
	"encoding/json"
	"fmt"
	"time"
)

type Collection struct {
	Name string 
	CreationTime time.Time
	Versioned bool
}

func ListCollections(requester Requester, credentials *Credentials) (
	[]Collection, error) {
	var response *Response
	var err error
	var result []Collection
	var rawSlice []map[string]interface{}

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

	err = json.Unmarshal(response.Body, &rawSlice)
    if err != nil {
        return nil, err
    }

    for _, rawMap := range rawSlice {
    	name := rawMap["name"].(string)
    	creationTime, err := time.Parse(time.RFC1123, 
    		rawMap["creation-time"].(string))
    	if err != nil {
    		return nil, err
    	}
    	versioning := rawMap["versioning"].(bool)
    	collection := Collection{name, creationTime, versioning}
    	result = append(result, collection)
    }

	return result, nil

}