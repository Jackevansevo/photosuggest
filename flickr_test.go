package main

import (
	"errors"
	"io/ioutil"
	"testing"
)

// [TODO] Figure out how to mock the getBytes function

func TestBuildFlickrURL(t *testing.T) {
	for _, c := range []struct {
		in, want string
		err      error
	}{
		{"dogs", "https://www.flickr.com/services/rest/?api_key=&format=json&media=photos&method=flickr.photos.search&nojsoncallback=1&per_page=30&safe_search=1&sort=relevance&text=dogs", nil},
		{"", "", errors.New("Specify a query string")},
	} {
		got, err := buildFlickURL(c.in)
		if c.err != nil && c.err.Error() != err.Error() {
			t.Errorf("buildFlickURL(%q) expected error: %q, got: %q", c.in, err, c.err)
		}
		if got != c.want {
			t.Errorf("buildFlickURL(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestProcessFlickrResponse(t *testing.T) {
	content, err := ioutil.ReadFile("fixtures/flickr.json")
	if err != nil {
		t.Errorf(err.Error())
	}

	out, err := processFlickrResponse(content)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := []flickrPhotoResult{
		{"https://farm6.staticflickr.com/5244/5340131446_3b7c380bea.jpg", "Flickr"},
		{"https://farm5.staticflickr.com/4026/4489119695_87144ba60b.jpg", "Flickr"},
		{"https://farm5.staticflickr.com/4131/4846208207_eb7d525741.jpg", "Flickr"},
	}

	for i := range expected {
		result := out[i].(flickrPhotoResult)
		if result != expected[i] {
			t.Errorf("Expected: %q, got: %q", result, expected[i])
		}
	}
}

func TestHandleMalformedJSON(t *testing.T) {
	_, err := processFlickrResponse([]byte(`[]`))
	expected := "json: cannot unmarshal array into Go value of type main.flickrAPIJsonResponse"
	if err == nil {
		t.Errorf("Expected error: %q, got: %T", expected, err)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected error: %q, got: %q", expected, err)
		}
	}
}
