package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var flickrURL, _ = url.Parse("https://www.flickr.com/services/rest/")

type flickrPhoto struct {
	ID     string
	Farm   int
	Server string
	Secret string
}

type flickrAPIJsonResponse struct {
	Data struct {
		PhotoList []flickrPhoto `json:"photo"`
	} `json:"photos"`
}

type flickrPhotoResult struct {
	URL    string `json:"url"`
	Source string `json:"source"`
}

var flickrParams = url.Values{
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

func queryFlickr(query string, client http.Client) ([]interface{}, error) {

	url, err := buildFlickURL(query)

	if err != nil {
		return nil, err
	}

	resp, err := getBytes(url, client)

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

	if query == "" {
		return "", errors.New("Specify a query string")
	}

	flickrParams.Set("text", query)
	flickrParams.Set("api_key", Env.FLICKR_API_KEY)

	flickrURL.RawQuery = flickrParams.Encode()

	return flickrURL.String(), nil
}

func processFlickrResponse(body []byte) ([]interface{}, error) {

	var resp flickrAPIJsonResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	results := make([]interface{}, len(resp.Data.PhotoList))

	for index, photo := range resp.Data.PhotoList {
		results[index] = flickrPhotoResult{photo.URL(), "Flickr"}
	}

	return results, nil
}
