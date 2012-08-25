package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultServiceDomain = "nimbus.io"
	defaultServicePort = 443
)

type client struct {
	credentials *Credentials
	httpClient  *http.Client
	serviceDomain string
	servicePort int
}

type Requester interface {
	DefaultHostName() string
	Request(method string, hostName string, path string) (*Response, error)
}

type Response struct {
	StatusCode int
	Status     string
	Body       []byte
}

func NewRequester(credentials *Credentials) (Requester, error) {
	serviceDomain := os.Getenv("NIMBUS_IO_SERVICE_DOMAIN")
	if serviceDomain == "" {
		serviceDomain = defaultServiceDomain
	}

	servicePortStr := os.Getenv("NIMBUS_IO_SERVICE_PORT")
	if servicePortStr == "" {
		servicePortStr = fmt.Sprintf("%d", defaultServicePort)
	}
	servicePort, err :=  strconv.Atoi(servicePortStr); if err != nil {
		return nil, err
	}

	requester := client{
		credentials,
		&http.Client{},
		serviceDomain,
		servicePort,
	}

	return &requester, nil
}

func (client *client) DefaultHostName() string {
	return fmt.Sprintf("%s:%d", client.serviceDomain, client.servicePort)
}

func (client *client) Request(method string, hostName string, path string) (
	*Response, error) {

	current_time := time.Now()
	timestamp := current_time.Unix()
	uri := fmt.Sprintf("http://%s%s", hostName, path)

	request, err := http.NewRequest(method, uri, nil); if err != nil {
		return nil, err
	}

	authString := ComputeAuthString(client.credentials, method, timestamp, path)
	request.Header.Add("Authorization", authString)
	request.Header.Add("x-nimbus-io-timestamp", fmt.Sprintf("%d", timestamp))
	request.Header.Add("agent", "gonimbusio/1.0")

	response, err := client.httpClient.Do(request); if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body); if err != nil {
		return nil, err
	}

	return &Response{response.StatusCode, response.Status, body}, nil
}
