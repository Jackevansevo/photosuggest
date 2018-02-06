package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const bingURL = "https://api.cognitive.microsoft.com/bing/v7.0/images/search"

var bingAPIKey = os.Getenv("BING_API_KEY")

type bingPhoto struct {
	URL string `json:"contentUrl"`
}

type bingJSON struct {
	PhotoList []bingPhoto `json:"value"`
}

func buildBingURL(query string) (string, error) {

	URL, err := url.Parse(bingURL)

	if err != nil {
		return "", nil
	}

	bingParams := url.Values{"q": {query}, "license": {"ModifyCommercially"}}

	URL.RawQuery = bingParams.Encode()

	return URL.String(), nil
}

func queryBing(query string) ([]string, error) {

	url, err := buildBingURL(query)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req.Header.Set("Ocp-Apim-Subscription-Key", bingAPIKey)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	urls, err := processBingResponse(body)

	if err != nil {
		return nil, err
	}

	return urls, nil

}

func processBingResponse(body []byte) ([]string, error) {

	var respJSON bingJSON

	if err := json.Unmarshal(body, &respJSON); err != nil {
		return nil, err
	}

	urls := make([]string, len(respJSON.PhotoList))

	for index, photo := range respJSON.PhotoList {
		urls[index] = photo.URL
	}

	return urls, nil

}
