package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
)

type Collection struct {
	Name string 
	CreationTime time.Time
	Versioned bool
}

const (
 	defaultCollectionPrefix = "dd"                                               
 	reservedCollectionPrefix = "rr"
)

func DefaultCollectionName(username string) string {                               
	return fmt.Sprintf("%s-%s", defaultCollectionPrefix, username) 
}

func ReservedCollectionName(username string, collectionName string) string {                
	return fmt.Sprintf("%s-%s-%s", reservedCollectionPrefix, username, 
 		collectionName)
 }

 func rawMapToCollection(rawMap map[string]interface{}) (*Collection, error) {
    name := rawMap["name"].(string)
    creationTime, err := time.Parse(time.RFC1123, 
    	rawMap["creation-time"].(string)); if err != nil {
    		return nil, err
    	}
    versioning := rawMap["versioning"].(bool)
    return &Collection{name, creationTime, versioning}, nil
}   
 
func ListCollections(requester Requester, credentials *Credentials) (
	[]Collection, error) {

	method := "GET"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections", credentials.Name)

	response, err := requester.Request(method, hostName, path); if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body); if err != nil {
		return nil, err
	}

	var rawSlice []map[string]interface{}
	err = json.Unmarshal(body, &rawSlice); if err != nil {
        return nil, err
    }

    var result []Collection
    for _, rawMap := range rawSlice {
    	collection, err := rawMapToCollection(rawMap); if err != nil {
    		return nil, err
    	}
    	result = append(result, *collection)
    }
	return result, nil
}

func CreateCollection(requester Requester, credentials *Credentials, 
	collectionName string) (*Collection, error) {

	method := "POST"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections?action=create&name=%s", 
		credentials.Name, collectionName)
	response, err := requester.Request(method, hostName, path); if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body); if err != nil {
		return nil, err
	}

	var rawMap map[string]interface{}
	err = json.Unmarshal(body, &rawMap); if err != nil {
        return nil, err
    }

    collection, err := rawMapToCollection(rawMap); if err != nil {
    	return nil, err
    }

	return collection, nil
}

func DeleteCollection(requester Requester, credentials *Credentials, 
	collectionName string) (bool, error) {

	method := "DELETE"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections/%s", 
		credentials.Name, collectionName)
	response, err := requester.Request(method, hostName, path); if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return false, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body); if err != nil {
		return false, err
	}

	var rawMap map[string]interface{}
	err = json.Unmarshal(body, &rawMap); if err != nil {
        return false, err
    }

    success := rawMap["success"].(bool)

	return success, nil
}
