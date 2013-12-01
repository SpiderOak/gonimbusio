package nimbusapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Collection struct {
	Name         string `json:"name"`
	CreationTime string `json:"creation-time"`
	Versioned    bool   `json:"versioning"`
}

type Key struct {
	Name              string `json:"key"`
	TimeStamp         string `json:"timestamp"`
	VersionIdentifier string `json:"version_identifier"`
}

type listKeyResult struct {
	Truncated bool  `json:"truncated"`
	KeySlice  []Key `json:"key_data"`
}

const (
	defaultCollectionPrefix  = "dd"
	reservedCollectionPrefix = "rr"
)

func DefaultCollectionName(username string) string {
	return fmt.Sprintf("%s-%s", defaultCollectionPrefix, username)
}

func ReservedCollectionName(username string, collectionName string) string {
	return fmt.Sprintf("%s-%s-%s", reservedCollectionPrefix, username,
		collectionName)
}

func ListCollections(requester Requester, userName string) (
	[]Collection, error) {

	method := "GET"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections", userName)

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return nil, err
	}

	response, err := requester.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result []Collection
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreateCollection(requester Requester, userName string,
	collectionName string) (*Collection, error) {

	method := "POST"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections?action=create&name=%s",
		userName, collectionName)

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return nil, err
	}

	response, err := requester.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var collection Collection
	err = json.Unmarshal(body, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func ListKeysInCollection(requester Requester, collectionName string) (
	[]Key, bool, error) {

	method := "GET"
	hostName := requester.CollectionHostName(collectionName)
	path := "/data/"

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return nil, false, err
	}

	response, err := requester.Do(request)
	if err != nil {
		return nil, false, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return nil, false, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, false, err
	}

	var listResult listKeyResult
	err = json.Unmarshal(body, &listResult)
	if err != nil {
		return nil, false, err
	}
	fmt.Printf("%v\n", listResult)

	return listResult.KeySlice, listResult.Truncated, nil
}

func ListVersionsInCollection(requester Requester, collectionName string, 
	prefix string) (
	[]Key, bool, error) {

	method := "GET"
	hostName := requester.CollectionHostName(collectionName)

	var path string
	if prefix == "" {
		path = "/?versions"
	} else {
		values := url.Values{}
		values.Set("prefix", prefix)
		path = fmt.Sprintf("/?versions&%s", values.Encode())
	}

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return nil, false, err
	}

	response, err := requester.Do(request)
	if err != nil {
		return nil, false, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return nil, false, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, false, err
	}

	var listResult listKeyResult
	err = json.Unmarshal(body, &listResult)
	if err != nil {
		return nil, false, err
	}
	fmt.Printf("%v\n", listResult)

	return listResult.KeySlice, listResult.Truncated, nil
}

func DeleteCollection(requester Requester, userName string,
	collectionName string) (bool, error) {

	method := "DELETE"
	hostName := requester.DefaultHostName()
	path := fmt.Sprintf("/customers/%s/collections/%s",
		userName, collectionName)

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return false, err
	}

	response, err := requester.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("DELETE %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return false, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	var resultMap map[string]bool
	err = json.Unmarshal(body, &resultMap)
	if err != nil {
		return false, err
	}

	return resultMap["success"], nil
}
