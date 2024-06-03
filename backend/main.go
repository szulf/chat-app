package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var dbs *dbService
var sessionKey []byte

func checkPassword(password string) bool {
	lowercase, err := regexp.MatchString("[a-z]", password)
	if err != nil {
		fmt.Println("error matching: ", err)
	}

	uppercase, err := regexp.MatchString("[A-Z]", password)
	if err != nil {
		fmt.Println("error matching: ", err)
	}

	digit, err := regexp.MatchString("[\\d]", password)
	if err != nil {
		fmt.Println("error matching: ", err)
	}

	special, err := regexp.MatchString("[#.?!@$%^&*]", password)
	if err != nil {
		fmt.Println("error matching: ", err)
	}

	length := len(password) > 7

	return uppercase && lowercase && digit && special && length
}

func checkCredentials(username, password, passwordConfirm string) (bool, string) {
	if len(username) < 3 {
		return false, "Username must be at least 4 characters long, and can't contain spaces"
	}

	if dbs.checkUsernameExists(username) {
		return false, "A user with that username already exists"
	}

	if !checkPassword(password) {
		return false, "Password needs 1 uppercase letter, 1 lowercase letter, 1 digit, 1 special character(#.?!@$%^&*), and needs to be at least 8 characters long"
	}

	if password != passwordConfirm {
		return false, "Passwords must match"
	}

	return true, ""
}

func writeSessionCookie(w http.ResponseWriter, sessionId string) error {
	authCookie := http.Cookie{
		Name:     "_session",
		Value:    sessionId,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
	}

	mac := hmac.New(sha256.New, sessionKey)
	mac.Write([]byte(authCookie.Name))
	mac.Write([]byte(authCookie.Value))
	signature := mac.Sum(nil)

	authCookie.Value = string(signature) + authCookie.Value
	authCookie.Value = base64.StdEncoding.EncodeToString([]byte(authCookie.Value))

	if len(authCookie.String()) > 4096 {
		return fmt.Errorf("cookie too long")
	}

	http.SetCookie(w, &authCookie)

	return nil
}

func read(r io.Reader) (string, error) {
	buffer := make([]byte, 4096)
	n, err := r.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(buffer[:n]), nil
}

type user struct {
	Username        string
	Password        string
	PasswordConfirm string
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("access-control-expose-headers", "Set-Cookie")
}

func login(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	switch r.Method {
	case "POST":
		requestBody, err := read(r.Body)
		if err != nil {
			fmt.Println("Error reading request body: ", err)
			os.Exit(1)
		}
		var u user
		json.Unmarshal([]byte(requestBody), &u)

		userId, err := dbs.getUserId(u.Username, u.Password)
		if err != nil {
			// Write user credientials wrong
			fmt.Println("Error getting userId: ", err)
			os.Exit(1)
		}

		sessionId := generateRandomString(128)
		err = writeSessionCookie(w, sessionId)
		if err != nil {
			fmt.Println("Error writing session cookie: ", err)
			os.Exit(1)
		}

		err = dbs.setSessionId(userId, sessionId)
		if err != nil {
			fmt.Println("Error setting sessionId: ", err)
			os.Exit(1)
		}

		content, err := os.ReadFile("htmls/chats.html")
		if err != nil {
			fmt.Println("Error reading from file: ", err)
		}
		w.Write(content)

	case "GET":
		content, err := os.ReadFile("htmls/login.html")
		if err != nil {
			fmt.Println("Error reading from file: ", err)
		}
		w.Write(content)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	switch r.Method {
	case "POST":
		requestBody, err := read(r.Body)
		if err != nil {
			fmt.Println("Error reading request body: ", err)
			os.Exit(1)
		}
		var u user
		json.Unmarshal([]byte(requestBody), &u)

		valid, msg := checkCredentials(u.Username, u.Password, u.PasswordConfirm)
		if !valid {
			w.Write([]byte("errMsg\n"))
			w.Write([]byte(msg))
			return
		}

		sessionId := generateRandomString(128)
		err = writeSessionCookie(w, sessionId)
		if err != nil {
			fmt.Println("Error writing session cookie: ", err)
			os.Exit(1)
		}

		userId, err := dbs.insertUser(u.Username, u.Password)
		if err != nil {
			fmt.Println("Error inserting user: ", err)
			os.Exit(1)
		}

		err = dbs.setSessionId(userId, sessionId)
		if err != nil {
			fmt.Println("Error setting sessionId: ", err)
			os.Exit(1)
		}

		content, err := os.ReadFile("htmls/chats.html")
		if err != nil {
			fmt.Println("Error reading from file: ", err)
		}
		w.Write(content)

	case "GET":
		content, err := os.ReadFile("htmls/register.html")
		if err != nil {
			fmt.Println("Error reading from file: ", err)
		}
		w.Write(content)
	}
}

func chats(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.Write([]byte("chats page"))
}

func main() {
	var err error
	sessionKey, err = hex.DecodeString(os.Getenv("CHATAPP_SESSION_KEY"))
	if err != nil {
		fmt.Println("Error decoding session key: ", err)
		os.Exit(1)
	}

	dbs, err = newDbConn(os.Getenv("CHATAPP_DB_CONN_URL"))
	if err != nil {
		fmt.Println("Error connecting to db: ", err)
		os.Exit(1)
	}
	defer dbs.closeConn()

	mux := http.NewServeMux()

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/chats", chats)

	log.Fatal(http.ListenAndServe("localhost:3000", mux))
}
