package nimbusapi

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

type Key struct {
	Name string 
	TimeStamp time.Time
	VersionIdentifier string
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
    	rawMap["creation-time"].(string))
    if err != nil {
    	return nil, err
    }
    versioning := rawMap["versioning"].(bool)
    return &Collection{name, creationTime, versioning}, nil
}   
 
func rawMapToKey(rawMap map[string]interface{}) (*Key, error) {
    name := rawMap["key"].(string)
    timeStamp, err := time.Parse(time.RFC1123, rawMap["timestamp"].(string))
    if err != nil {
    	return nil, err
    }
    versionIdentifier := rawMap["version_identifier"].(string)
    return &Key{name, timeStamp, versionIdentifier}, nil
}   
 
func ListCollections(requester Requester, userName string) (
	[]Collection, error) {

	method := "GET"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections", userName)

	response, err := requester.Request(method, hostName, path, nil)
	if err != nil {
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

func CreateCollection(requester Requester, userName string, 
	collectionName string) (*Collection, error) {

	method := "POST"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections?action=create&name=%s", 
		userName, collectionName)
	response, err := requester.Request(method, hostName, path, nil)
	if err != nil {
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

func ListKeysInCollection(requester Requester, collectionName string) (
	[]Key, bool, error) {

	method := "GET"
	hostName := requester.CollectionHostName(collectionName)
	path := "/data/"

	response, err := requester.Request(method, hostName, path, nil) 
	if err != nil {
		return nil, false, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return nil, false, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body); if err != nil {
		return nil, false, err
	}

	var rawResultMap map[string]interface{}
	err = json.Unmarshal(body, &rawResultMap); if err != nil {
        return nil, false, err
    }
    fmt.Printf("%v\n", rawResultMap)

	rawKeySlice, ok := rawResultMap["key_data"].([]interface{})
	if !ok {
		err = fmt.Errorf("Unable to convert %v", rawResultMap["key_data"])
        return nil, false, err
    }

    var keySlice []Key
    for _, rawKeyInterface := range rawKeySlice {
		rawKeyMap, ok := rawKeyInterface.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("Unable to convert %v", rawKeyInterface)
        	return nil, false, err
    	}
    	key, err := rawMapToKey(rawKeyMap); if err != nil {
    		return nil, false, err
    	}
    	keySlice = append(keySlice, *key)
    }

    truncated := rawResultMap["truncated"].(bool)

 	return keySlice, truncated, nil
}

func DeleteCollection(requester Requester, userName string, 
	collectionName string) (bool, error) {

	method := "DELETE"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections/%s", 
		userName, collectionName)
	response, err := requester.Request(method, hostName, path, nil) 
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("DELETE %s %s failed (%d) %s", hostName, path, 
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
