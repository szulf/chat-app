package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html")
		x := make([]byte, 4096)
		n, err := r.Body.Read(x)
		if err != nil && err != io.EOF {
			fmt.Println("what?", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(x[:n]))

		_, err = w.Write([]byte("This got fetched!<br>Yupee!"))
		if err != nil {
			fmt.Println("what the?", err.Error())
			os.Exit(1)
		}
	})

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
