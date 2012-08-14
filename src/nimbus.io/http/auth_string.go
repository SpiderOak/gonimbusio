package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func ComputeAuthString(credentials *Credentials, method string, timestamp int64,
	uri string) string {

	message := fmt.Sprintf("%s\n%s\n%d\n%s", credentials.name, method,
		timestamp, uri)
	h := hmac.New(sha256.New, credentials.authKey)
	h.Write([]byte(message))
	return fmt.Sprintf("NIMBUS.IO %d:%x", credentials.authKeyId, h.Sum(nil))
}
