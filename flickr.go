package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

const flickrURL = "https://www.flickr.com/services/rest/"

type flickrPhoto struct {
	ID     string
	Farm   int
	Server string
	Secret string
}

type flickJSON struct {
	Data struct {
		PhotoList []flickrPhoto `json:"photo"`
	} `json:"photos"`
}

var flickrAPIKey = os.Getenv("FLICKR_API_KEY")

var flickrParams = url.Values{
	"api_key":        {flickrAPIKey},
	"format":         {"json"},
	"media":          {"photos"},
	"method":         {"flickr.photos.search"},
	"nojsoncallback": {"1"},
	"per_page":       {"30"},
	"safe_search":    {"1"},
	"sort":           {"relevance"},
}

func (p flickrPhoto) URL() string {
	urlFmt := "https://farm%d.staticflickr.com/%s/%s_%s.jpg"
	return fmt.Sprintf(urlFmt, p.Farm, p.Server, p.ID, p.Secret)
}

func queryFlickr(query string) ([]string, error) {

	url, err := buildFlickURL(query)

	if err != nil {
		return nil, err
	}

	resp, err := HTTPGetBytes(url)

	if err != nil {
		return nil, err
	}

	json, err := processFlickrResponse(resp)

	if err != nil {
		return nil, err
	}

	return json, nil
}

func buildFlickURL(query string) (string, error) {
	URL, err := url.Parse(flickrURL)

	if err != nil {
		return "", err
	}

	flickrParams.Set("text", query)

	URL.RawQuery = flickrParams.Encode()

	return URL.String(), nil
}

func processFlickrResponse(body []byte) ([]string, error) {

	var respJSON flickJSON

	if err := json.Unmarshal(body, &respJSON); err != nil {
		return nil, err
	}

	urls := make([]string, len(respJSON.Data.PhotoList))

	for index, photo := range respJSON.Data.PhotoList {
		urls[index] = photo.URL()
	}

	return urls, nil
}
