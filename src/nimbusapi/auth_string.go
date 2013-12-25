/*Package nimbusapi provides go routines to access the nimbus.io REST API
 */
package nimbusapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/url"
)

// ComputeAuthString returns the authentication string sent to nimbus.io
// web servers
func ComputeAuthString(credentials *Credentials, method string, timestamp int64,
	path string) (string, error) {

	rawPath, err := url.QueryUnescape(path)
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf("%s\n%s\n%d\n%s", credentials.Name, method,
		timestamp, rawPath)
	h := hmac.New(sha256.New, credentials.AuthKey)
	h.Write([]byte(message))
	authString := fmt.Sprintf("NIMBUS.IO %d:%x", credentials.AuthKeyID,
		h.Sum(nil))
	return authString, nil
}
