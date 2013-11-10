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

type Credentials struct {
	Name      string
	AuthKeyId int
	AuthKey   []byte
}

const credentialsFileName = ".nimbus.io"

func (credentials *Credentials) Equal(other *Credentials) bool {
	return credentials.Name == other.Name &&
		credentials.AuthKeyId == other.AuthKeyId &&
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
		return nil, errors.New("can't parse AuthKeyId")
	}
	credentials.AuthKeyId, err = strconv.Atoi(fields[1])
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

func LoadCredentialsFromPath(path string) (*Credentials, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return loadCredentials(file)
}

func LoadCredentialsFromDefault() (*Credentials, error) {
	userRec, err := user.Current()
	if err != nil {
		return nil, err
	}
	credentialsPath := path.Join(userRec.HomeDir, credentialsFileName)
	return LoadCredentialsFromPath(credentialsPath)
}
