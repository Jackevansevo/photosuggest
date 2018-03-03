package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var bingURL, _ = url.Parse("https://api.cognitive.microsoft.com/bing/v7.0/images/search")

type bingAPIPhoto struct {
	URL string `json:"contentUrl"`
}

type bingAPIJsonResponse struct {
	PhotoList []bingAPIPhoto `json:"value"`
}

type bingPhotosuggestResult struct {
	URL    string `json:"url"`
	Source string `json:"source"`
}

func buildBingURL(query string) (string, error) {

	if query == "" {
		return "", errors.New("Specify a query string")
	}

	bingParams := url.Values{"q": {query}, "license": {"ModifyCommercially"}}

	bingURL.RawQuery = bingParams.Encode()

	return bingURL.String(), nil
}

func queryBing(query string, client http.Client) ([]interface{}, error) {

	url, err := buildBingURL(query)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", Env.BingAPIKey)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	// Report any server errors
	if resp.StatusCode >= 400 {
		return nil, errors.New("Bing: " + string(body))
	}

	urls, err := processBingResponse(body)

	if err != nil {
		return nil, err
	}

	return urls, nil

}

func processBingResponse(body []byte) ([]interface{}, error) {

	var resp bingAPIJsonResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	results := make([]interface{}, len(resp.PhotoList))

	for index, photo := range resp.PhotoList {
		results[index] = bingPhotosuggestResult{photo.URL, "Bing"}
	}

	return results, nil

}
