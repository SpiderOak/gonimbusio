package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultServiceDomain = "nimbus.io"
	defaultServicePort   = 443
)

type client struct {
	credentials   *Credentials
	httpClient    *http.Client
	serviceDomain string
	servicePort   int
}

type Requester interface {
	DefaultHostName() string
	CollectionHostName(collectionName string) string
	Request(method string, hostName string, path string, body io.Reader) (
		*Response, error)
}

type Response struct {
	StatusCode int
	Status     string
	Body       io.ReadCloser
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
	servicePort, err := strconv.Atoi(servicePortStr)
	if err != nil {
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

func (client *client) CollectionHostName(collectionName string) string {
	return fmt.Sprintf("%s.%s", collectionName, client.DefaultHostName())
}

func (client *client) Request(method string, hostName string, path string,
	body io.Reader) (*Response, error) {

	current_time := time.Now()
	timestamp := current_time.Unix()
	uri := fmt.Sprintf("http://%s%s", hostName, path)

	request, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	authString := ComputeAuthString(client.credentials, method, timestamp, path)
	request.Header.Add("Authorization", authString)
	request.Header.Add("x-nimbus-io-timestamp", fmt.Sprintf("%d", timestamp))
	request.Header.Add("agent", "gonimbusio/1.0")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return &Response{response.StatusCode, response.Status, response.Body}, nil
}
