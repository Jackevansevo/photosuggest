package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// [TODO] Test without internet

func TestMissingAPIKeys(t *testing.T) {

	// Take snapshot of the environment
	env := os.Environ()

	// Temporarily clear all environment variables
	os.Clearenv()

	defer func() {
		if r := recover(); r != nil {

			expectedErr := "required key BING_API_KEY missing value"

			// Cast recovery type to error type
			err, ok := r.(error)
			if !ok {
				t.Errorf("pkg: %v", r)
			}

			if err.Error() != expectedErr {
				t.Errorf("Expected: %v, got: %v", expectedErr, err)
			}
		}
	}()

	main()

	// Recover environment
	for _, pair := range env {
		splits := strings.Split(pair, "=")
		key, val := splits[0], splits[1]
		err := os.Setenv(key, val)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

}

func TestQueryWithMissingParam(t *testing.T) {

	server := httptest.NewServer(handler())
	defer server.Close()

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	expectedErr := "Specify query"

	if strings.TrimSpace(string(b)) != expectedErr {
		t.Errorf("expected body to contain: %v, got: %v", expectedErr, string(b))
	}

}
