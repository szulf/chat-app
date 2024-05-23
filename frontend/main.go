package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func setContentType(filePath string, writer http.ResponseWriter) {
	splitStrings := strings.Split(filePath, ".")
	extension := splitStrings[len(splitStrings)-1]

	switch extension {
	case "html":
		writer.Header().Set("Content-Type", "text/html")

	case "js":
		writer.Header().Set("Content-Type", "text/javascript")

	case "css":
		writer.Header().Set("Content-Type", "text/css")

	default:
		writer.Header().Set("Content-Type", "text/plain")
	}
}

func notFoundSite(w http.ResponseWriter, r *http.Request) {
	file := "err404.html"
	content, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			notFoundSite(w, r)
		} else {
			fmt.Println("Error reading from file: ", err.Error())
		}
	}
	setContentType(file, w)

	_, err = w.Write(content)
	if err != nil {
		fmt.Println("Error writing to http: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile(r.RequestURI[1:])
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("not found error occurred")
				notFoundSite(w, r)
			} else {
				fmt.Println("Error reading from file: ", err.Error())
			}
		}
		setContentType(r.RequestURI[1:], w)

		_, err = w.Write(content)
		if err != nil {
			fmt.Println("Error writing to http: ", err.Error())
			os.Exit(1)
		}
	})

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
