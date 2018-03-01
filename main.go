package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var client = http.Client{Timeout: time.Duration(5 * time.Second)}

type EnvSpec struct {
	BING_API_KEY   string `envconfig:"BING_API_KEY" required:"true"`
	FLICKR_API_KEY string `envconfig:"FLICKR_API_KEY" required:"true"`
}

var Env EnvSpec

func indexHandler(w http.ResponseWriter, r *http.Request) {

	// Make Channels
	resc, errc := make(chan []interface{}), make(chan error)

	query := r.URL.Query().Get("q")

	queryFuncs := []func(string, http.Client) ([]interface{}, error){
		queryBing,
		queryFlickr,
	}

	if query == "" {
		http.Error(w, "Specify query", http.StatusBadRequest)
		return
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
	w.Write(resp)

}

func main() {

	// Loads API keys from environment
	envconfig.MustProcess("", &Env)

	log.Fatal(http.ListenAndServe(":3001", handler()))
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	return mux
}
