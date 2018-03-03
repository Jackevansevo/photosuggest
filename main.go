package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var client = http.Client{Timeout: time.Duration(5 * time.Second)}

type envSpec struct {
	BingAPIKey   string `envconfig:"BING_API_KEY" required:"true"`
	FlickrAPIKey string `envconfig:"FLICKR_API_KEY" required:"true"`
}

// Env contains environment variables
var Env envSpec

var sourceFuncs = map[string]func(string, http.Client) ([]interface{}, error){
	"bing":   queryBing,
	"flickr": queryFlickr,
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	// Make Channels
	resc, errc := make(chan []interface{}), make(chan error)

	q := r.URL.Query()

	query := q.Get("q")

	if query == "" {
		http.Error(w, "Specify query", http.StatusBadRequest)
		return
	}

	// Parse the sources from the url query params
	sources := q.Get("sources")
	queryFuncs := make([]func(string, http.Client) ([]interface{}, error), 0)

	if sources != "" {

		for _, source := range strings.Split(sources, " ") {
			fn, ok := sourceFuncs[source]
			if ok {
				queryFuncs = append(queryFuncs, fn)
			}
		}
	} else {
		// Use all sources by default
		for _, fn := range sourceFuncs {
			queryFuncs = append(queryFuncs, fn)
		}
	}

	for _, fn := range queryFuncs {

		go func(query string, resc chan []interface{},
			errc chan error, fn func(string, http.Client) ([]interface{}, error)) {
			results, err := fn(query, client)
			if err != nil {
				errc <- err
				return
			}
			resc <- results
		}(query, resc, errc, fn)

	}

	results := make([]interface{}, 0)

	for _ = range queryFuncs {
		select {
		case res := <-resc:
			results = append(results, res...)
		case err := <-errc:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(results)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {

	// Loads API keys from environment
	envconfig.MustProcess("", &Env)

	log.Fatal(http.ListenAndServe(":8000", handler()))
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	return mux
}
