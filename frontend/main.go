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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file := ""
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			file = "index.html"
		} else {
			file = r.RequestURI[1:]
		}

		content, err := os.ReadFile(file)
		if err != nil {
			if os.IsNotExist(err) {
				// Error 404 implementation here
				fmt.Println("not found error")
			}
			fmt.Println("Error reading from file: ", err.Error())
		}

		setContentType(file, w)
		_, err = w.Write(content)
		if err != nil {
			fmt.Println("Error writing to http: ", err.Error())
			os.Exit(1)
		}
	})

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
