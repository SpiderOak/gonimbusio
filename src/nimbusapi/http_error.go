/* package nimbusapi provides go reioutines to access the nimbus.io REST API
 */
package nimbusapi

import (
	"fmt"
)

type HTTPError struct {
	StatusCode int
	Message    string
}

func (httpError HTTPError) Error() string {
	return fmt.Sprintf("(%d) %s", httpError.StatusCode, httpError.Message)
}
