package nimbusapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"io"
	"io/ioutil"
	"strconv"
)

func StartConjoined(requester Requester, collectionName string, key string) (
	string, error) {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "start")
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key), 
		values.Encode())

	response, err := requester.Request(method, hostName, path, nil) 
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
	var conjoinedMap map[string]string
	err = json.Unmarshal(responseBody, &conjoinedMap); if err != nil {
        return "", err
    }

    return conjoinedMap["conjoined_identifier"], nil
}

func AbortConjoined(requester Requester, collectionName string, key string, 
	conjoinedIdentifier string) error {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "abort")
	values.Set("conjoined_identifier", conjoinedIdentifier)
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key), 
		values.Encode())

	response, err := requester.Request(method, hostName, path, nil) 
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body); if err != nil {
		return err
	}
	var conjoinedMap map[string]string
	err = json.Unmarshal(responseBody, &conjoinedMap); if err != nil {
        return err
    }

    if conjoinedMap["success"] != "true" {
    	return fmt.Errorf("conjoined action=abort returned %s", 
    		conjoinedMap["success"])
    }

    return nil
}
func FinishConjoined(requester Requester, collectionName string, key string, 
	conjoinedIdentifier string) error {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "finish")
	values.Set("conjoined_identifier", conjoinedIdentifier)
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key), 
		values.Encode())

	response, err := requester.Request(method, hostName, path, nil) 
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body); if err != nil {
		return err
	}
	var conjoinedMap map[string]string
	err = json.Unmarshal(responseBody, &conjoinedMap); if err != nil {
        return err
    }

    if conjoinedMap["success"] != "true" {
    	return fmt.Errorf("conjoined action=abort returned %s", 
    		conjoinedMap["success"])
    }

    return nil
}

type ConjoinedParams struct {
	ConjoinedIdentifier string 
	ConjoinedPart int
}

func Archive(requester Requester, collectionName string, key string, 
	conjoinedParams *ConjoinedParams, requestBody io.Reader) (string, error) {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)

	var path string
	if conjoinedParams != nil {
		values := url.Values{}
		values.Set("conjoined_identifier", conjoinedParams.ConjoinedIdentifier)
		values.Set("conjoined_part", 
			strconv.Itoa(conjoinedParams.ConjoinedPart))
		path = fmt.Sprintf("/data/%s?%s", url.QueryEscape(key), values.Encode())
	} else {
		path = fmt.Sprintf("/data/%s", url.QueryEscape(key))
	}

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
	var versionMap map[string]string
	err = json.Unmarshal(responseBody, &versionMap); if err != nil {
        return "", err
    }

    return versionMap["version_identifier"], nil
}

func Retrieve(requester Requester, collectionName string, key string) (
	io.ReadCloser, error) {
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

func DeleteKey(requester Requester, collectionName string, key string) (error) {
	method := "DELETE"
	hostName := requester.CollectionHostName(collectionName)
	path := fmt.Sprintf("/data/%s", url.QueryEscape(key))

	response, err := requester.Request(method, hostName, path, nil) 
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path, 
			response.StatusCode, response.Body)
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body); if err != nil {
		return err
	}

	var resultMap map[string]bool
	err = json.Unmarshal(responseBody, &resultMap); if err != nil {
        return err
    }
    if !resultMap["success"] {
    	err = fmt.Errorf("unexpected 'false' for 'success'")
    	return err
    }

	return nil
}
