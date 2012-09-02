package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"io"
	"io/ioutil"
)

func Archive(requester Requester, credentials *Credentials, 
	collectionName string, key string, requestBody io.Reader) (string, error) {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	path := fmt.Sprintf("/data/%s", url.QueryEscape(key))

	response, err := requester.Request(method, hostName, path, requestBody) 
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return "", err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body); if err != nil {
		return "", err
	}
	var rawMap map[string]interface{}
	err = json.Unmarshal(responseBody, &rawMap); if err != nil {
        return "", err
    }
    versionIdentifier := rawMap["version_identifier"].(string)

	return versionIdentifier, nil
}

func Retrieve(requester Requester, credentials *Credentials, 
	collectionName string, key string) (io.ReadCloser, error) {
	method := "GET"
	hostName := requester.CollectionHostName(collectionName)
	path := fmt.Sprintf("/data/%s", url.QueryEscape(key))

	response, err := requester.Request(method, hostName, path, nil) 
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return nil, err
	}

	return response.Body, nil
}
