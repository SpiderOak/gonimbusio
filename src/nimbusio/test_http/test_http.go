package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	nimbusiohttp "nimbusio/http"
)

func main() {
	fmt.Println("start")
	credentials := nimbusiohttp.Credentials{}
	method := "GET"
	current_time := time.Now()
	timestamp := current_time.Unix()
	uri := "http://dev.nimbus.io:9000/customers/xxx/collections"

	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalf("NewRequest failed %s\n", err)
	}
	fmt.Printf("req = %v\n", req)

	authString := nimbusiohttp.ComputeAuthString(&credentials, method, timestamp, 
		uri)
	req.Header.Add("Authorization", authString)
	req.Header.Add("x-nimbus-io-timestamp", string(timestamp))
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
