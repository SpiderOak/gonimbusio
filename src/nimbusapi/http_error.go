/*Package nimbusapi provides go routines to access the nimbus.io REST API
 */
package nimbusapi

import (
	"fmt"
)

// HTTPError implements the error interface with the additional information
// of the specific HTTP status that is in error
type HTTPError struct {
	StatusCode int
	Message    string
}

// HTTPError.Error implments th error interface by returning a string
func (httpError HTTPError) Error() string {
	return fmt.Sprintf("(%d) %s", httpError.StatusCode, httpError.Message)
}
