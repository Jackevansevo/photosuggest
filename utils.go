package main

import (
	"io/ioutil"
	"net/http"
)

// HTTPGetBytes Makes a Http Get request and returns body as bytes
func HTTPGetBytes(url string) ([]byte, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil

}

func concatLists(items ...[]string) (out []string) {
	for _, l := range items {
		out = append(out, l...)
	}
	return out
}
