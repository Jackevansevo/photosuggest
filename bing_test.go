package main

import (
	"errors"
	"io/ioutil"
	"testing"
)

func TestBuildBingURL(t *testing.T) {
	for _, c := range []struct {
		in, want string
		err      error
	}{
		{"dogs", "https://api.cognitive.microsoft.com/bing/v7.0/images/search?license=ModifyCommercially&q=dogs", nil},
		{"hello world", "https://api.cognitive.microsoft.com/bing/v7.0/images/search?license=ModifyCommercially&q=hello+world", nil},
		{"", "", errors.New("Specify a query string")},
	} {
		got, err := buildBingURL(c.in)
		if c.err != nil && c.err.Error() != err.Error() {
			t.Errorf("buildBingURL(%q) expected error: %q, got: %q", c.in, err, c.err)
		}
		if got != c.want {
			t.Errorf("buildBingURL(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestProcessBingResponse(t *testing.T) {

	content, err := ioutil.ReadFile("fixtures/bing.json")
	if err != nil {
		t.Errorf(err.Error())
	}

	out, err := processBingResponse(content)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := []bingPhotosuggestResult{
		{"http://upload.wikimedia.org/wikipedia/commons/a/a2/Flickr_-_ggallice_-_Street_dogs_%281%29.jpg", "Bing"},
		{"http://1.bp.blogspot.com/-eO4Bzg4yAuA/TiRKxZTB6VI/AAAAAAAADkQ/8rDy98dId8w/s1600/dogs.jpg", "Bing"},
		{"http://upload.wikimedia.org/wikipedia/commons/4/4d/Stray_dogs_crosswalk.jpg", "Bing"},
	}

	for i := range expected {
		result := out[i].(bingPhotosuggestResult)
		if result != expected[i] {
			t.Errorf("Expected: %q, got: %q", result, expected[i])
		}
	}
}

func TestProcessBingResponseWithMalformedResponse(t *testing.T) {
	_, err := processBingResponse([]byte(`[]`))
	expected := "json: cannot unmarshal array into Go value of type main.bingAPIJsonResponse"
	if err == nil {
		t.Errorf("Expected error: %q, got: %T", expected, err)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected error: %q, got: %q", expected, err)
		}
	}
}
