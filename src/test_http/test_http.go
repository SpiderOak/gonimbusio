package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Credentials struct {
	name      string
	authKeyId int
	authKey   []byte
}

func ComputeAuthString(credentials *Credentials, method string, timestamp int64,
	uri string) string {

	message := fmt.Sprintf("%s\n%s\n%d\n%s\n", credentials.name, method,
		timestamp, uri)
	h := hmac.New(sha256.New, credentials.authKey)
	h.Write([]byte(message))
	return fmt.Sprintf("NIMBUS.IO %d:%x", credentials.authKeyId, h.Sum(nil))
}

func main() {
	fmt.Println("start")
	credentials := Credentials{}
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

	authString := ComputeAuthString(&credentials, method, timestamp, uri)
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
