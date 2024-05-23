package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Write([]byte("did this get fetched?"))
	})

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
