package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	nimbusiohttp "nimbus.io/http"
	"time"
)

const (
	base_address = "dev.nimbus.io:9000"
)

func main() {
	fmt.Println("start")
	var credentials *nimbusiohttp.Credentials
	var err error

	sp := flag.String("credentials", "", "path to credentials file")
	flag.Parse()
	if *sp == "" {
		credentials, err = nimbusiohttp.LoadCredentialsFromDefault()
	} else {
		credentials, err = nimbusiohttp.LoadCredentialsFromPath(*sp)
	}
	if err != nil {
		log.Fatalf("Error loading credentials %s\n", err)
	}

	method := "GET"
	current_time := time.Now()
	timestamp := current_time.Unix()
	base_uri := fmt.Sprintf("/customers/%s/collections", credentials.Name)
	uri := fmt.Sprintf("http://%s%s", base_address, base_uri)

	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalf("NewRequest failed %s\n", err)
	}
	fmt.Printf("req = %v\n", req)

	authString := nimbusiohttp.ComputeAuthString(credentials, method, timestamp,
		base_uri)
	req.Header.Add("Authorization", authString)
	req.Header.Add("x-nimbus-io-timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Add("agent", "gonimbusio/1.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("DO failed %s\n", err)
	}
	defer resp.Body.Close()
	fmt.Printf("resp = %v\n", resp)

	//body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("end")
}
