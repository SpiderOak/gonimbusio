/*Package nimbusapi provides go routines to access the nimbus.io REST API
 */
package nimbusapi

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
)

// Credentials are a user's nimbus.io credentials
type Credentials struct {
	Name      string
	AuthKeyID int
	AuthKey   []byte
}

const credentialsFileName = ".nimbus.io"

// Equal returns true if two sets of credentials are the same
func (credentials *Credentials) Equal(other *Credentials) bool {
	return credentials.Name == other.Name &&
		credentials.AuthKeyID == other.AuthKeyID &&
		bytes.Equal(credentials.AuthKey, other.AuthKey)
}

func loadCredentials(reader io.Reader) (*Credentials, error) {
	credentials := Credentials{}
	bufferedReader := bufio.NewReader(reader)

	line, err := bufferedReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(line)
	if len(fields) != 2 || fields[0] != "Username" {
		return nil, errors.New("can't parse Username")
	}
	credentials.Name = fields[1]

	line, err = bufferedReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fields = strings.Fields(line)
	if len(fields) != 2 || fields[0] != "AuthKeyId" {
		return nil, errors.New("can't parse AuthKeyID")
	}
	credentials.AuthKeyID, err = strconv.Atoi(fields[1])
	if err != nil {
		return nil, err
	}

	line, err = bufferedReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fields = strings.Fields(line)
	if len(fields) != 2 || fields[0] != "AuthKey" {
		return nil, errors.New("can't parse AuthKey")
	}
	credentials.AuthKey = []byte(fields[1])

	return &credentials, nil
}

// LoadCredentialsFromPath loads Credentials from a file at a specified path
func LoadCredentialsFromPath(path string) (*Credentials, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return loadCredentials(file)
}

// LoadCredentialsFromDefault loads Credentials from a file at a default
// location
func LoadCredentialsFromDefault() (*Credentials, error) {
	userRec, err := user.Current()
	if err != nil {
		return nil, err
	}
	credentialsPath := path.Join(userRec.HomeDir, credentialsFileName)
	return LoadCredentialsFromPath(credentialsPath)
}
