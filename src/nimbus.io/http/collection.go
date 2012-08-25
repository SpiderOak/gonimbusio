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

	method := "GET"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections", credentials.Name)

	response, err := requester.Request(method, hostName, path); if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return nil, err
	}

	var rawSlice []map[string]interface{}
	err = json.Unmarshal(response.Body, &rawSlice); if err != nil {
        return nil, err
    }

    var result []Collection
    for _, rawMap := range rawSlice {
    	name := rawMap["name"].(string)
    	creationTime, err := time.Parse(time.RFC1123, 
    		rawMap["creation-time"].(string)); if err != nil {
    		return nil, err
    	}
    	versioning := rawMap["versioning"].(bool)
    	collection := Collection{name, creationTime, versioning}
    	result = append(result, collection)
    }

	return result, nil

}

