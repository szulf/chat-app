package main

import (
	"log"
	"net/http"
)

type requests map[string]func(w http.ResponseWriter, r *http.Request) error

func home(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("this is in the main"))
	return nil
}

func errSite(w http.ResponseWriter) {
	w.Write([]byte("This doesn't exist!"))
}

func main() {
	rqs := make(requests)
	rqs["home"] = home
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		if r.RequestURI == "/" {
			rqs["home"](w, r)
			return
		}

		fn, exists := rqs[r.RequestURI[1:]]
		if !exists {
			errSite(w)
			return
		}
		fn(w, r)
	})

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
