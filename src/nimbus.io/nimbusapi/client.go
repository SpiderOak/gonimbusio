package nimbusapi

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
	protocol string
	httpClient    *http.Client
	serviceDomain string
	servicePort   int
}

type Requester interface {
	DefaultHostName() string
	CollectionHostName(collectionName string) string
	CreateRequest(method string, hostName string, path string, 
		body io.Reader) (*http.Request, error)
	Do(*http.Request) (*http.Response, error)
}

func NewRequester(credentials *Credentials) (Requester, error) {
	serviceDomain := os.Getenv("NIMBUS_IO_SERVICE_DOMAIN")
	if serviceDomain == "" {
		serviceDomain = defaultServiceDomain
	}

	var protocol string
	useSSL := os.Getenv("NIMBUS_IO_SERVICE_SSL")
	if useSSL == "" || useSSL == "1" {
		protocol = "https"
	} else {
		protocol = "http"
	}

	servicePortStr := os.Getenv("NIMBUS_IO_SERVICE_PORT")
	if servicePortStr == "" {
		servicePortStr = fmt.Sprintf("%d", defaultServicePort)
	}
	servicePort, err := strconv.Atoi(servicePortStr)
	if err != nil {
		return nil, err
	}

	// TODO: the web server is sending a "Connection: close" header
	// We should deal with that
	httpTransport := &http.Transport{DisableKeepAlives: true}
	httpClient := &http.Client{Transport: httpTransport}

	requester := client{
		credentials,
		protocol,
		httpClient,
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

// Create a http.Request with our custom authentication headers
func (client *client) CreateRequest(method string, hostName string, path string,
	body io.Reader) (*http.Request, error) {

	current_time := time.Now()
	timestamp := current_time.Unix()
	uri := fmt.Sprintf("%s://%s%s", client.protocol, hostName, path)

	request, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	authString, err := ComputeAuthString(client.credentials, method, timestamp,
		path)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", authString)
	request.Header.Add("x-nimbus-io-timestamp", fmt.Sprintf("%d", timestamp))
	request.Header.Add("agent", "gonimbusio/1.0")

	return request, nil
}

func (client *client) Do(request *http.Request) (*http.Response, error) {
	return client.httpClient.Do(request)
}
