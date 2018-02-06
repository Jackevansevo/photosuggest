package main

import (
	"encoding/json"
	"net/http"

	"github.com/rs/cors"
)

func fetchResults(query string, resc chan []string, errc chan error, fn func(string) ([]string, error)) {
	urls, err := fn(query)
	if err != nil {
		errc <- err
		return
	}
	resc <- urls
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Make URL Channel
	resc, errc := make(chan []string), make(chan error)

	queryFuncs := []func(string) ([]string, error){
		queryBing,
		queryFlickr,
	}

	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "Specify query", http.StatusBadRequest)
		return
	}

	for _, fn := range queryFuncs {
		go fetchResults(query, resc, errc, fn)
	}

	urls := []string{}

	for i := 0; i < len(queryFuncs); i++ {
		select {
		case res := <-resc:
			urls = append(urls, res...)
		case err := <-errc:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	js, err := json.Marshal(urls)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":3001", handler)
}
