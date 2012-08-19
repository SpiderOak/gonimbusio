package http

import (
	"strings"
	"testing"
)

type testCredentialsEntry struct {
	testString  string
	credentials Credentials
}

// test data generated from python library lumberyard
var testCredentialsData = []testCredentialsEntry{
	testCredentialsEntry{
		`Username motoboto-benchmark-000
AuthKeyId 1
AuthKey 4TVjaSkh5GNENBqs+GX2OUrDQlofOgzP/0QB+F1+TYY
`,		
		Credentials{
			"motoboto-benchmark-000",
			1,
			[]byte("4TVjaSkh5GNENBqs+GX2OUrDQlofOgzP/0QB+F1+TYY"),
		},
	},
}

func TestLoadCredentials(t *testing.T) {
	for _, entry := range testCredentialsData {
		reader := strings.NewReader(entry.testString)
		credentials, err := loadCredentials(reader)
		if err != nil {
			t.Fatalf("error %s, %v", err, entry)
		}
		if !credentials.Equal(&entry.credentials) {
			t.Fatalf("credentials mismatch %v, %v", credentials, entry)
		}
	}
}
