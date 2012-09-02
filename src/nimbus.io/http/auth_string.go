package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/url"

)

func ComputeAuthString(credentials *Credentials, method string, timestamp int64,
	path string) (string, error) {

	rawPath, err := url.QueryUnescape(path); if err != nil {
		return "", err
	}

	message := fmt.Sprintf("%s\n%s\n%d\n%s", credentials.Name, method,
		timestamp, rawPath)
	h := hmac.New(sha256.New, credentials.AuthKey)
	h.Write([]byte(message))
	authString := fmt.Sprintf("NIMBUS.IO %d:%x", credentials.AuthKeyId, 
		h.Sum(nil))
	return authString, nil
}
