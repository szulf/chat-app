package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type user struct {
	Username string
	Password string
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("access-control-expose-headers", "Set-Cookie")
}

func login(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	if r.Method == "POST" {
		buffer := make([]byte, 4096)
		n, err := r.Body.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("error reading: ", err)
			os.Exit(1)
		}
		var u user
		json.Unmarshal(buffer[:n], &u)

		authCookie := http.Cookie{
			Name:  "test",
			Value: base64.StdEncoding.EncodeToString([]byte("test")),
		}

		// Have to set cookie before writing to response
		http.SetCookie(w, &authCookie)
		w.Write([]byte("Sent cookie"))
		w.Write([]byte("\n"))

		w.Write([]byte(u.Username))
		w.Write([]byte("\n"))
		w.Write([]byte(u.Password))
		w.Write([]byte("\n"))
	}

	w.Write([]byte(r.Method))
	w.Write([]byte("\n"))
	w.Write([]byte("login page"))
	w.Write([]byte("\n"))

}

func register(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.Write([]byte("register page"))
}

func chats(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.Write([]byte("chats page"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/chats", chats)

	log.Fatal(http.ListenAndServe("localhost:3000", mux))
}
