package main

import (
	"database/sql"
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(n int) string {
	userId := ""
	for i := 0; i < n; i++ {
		userId += string(chars[rand.Intn(len(chars))])
	}
	return userId
}

type dbService struct {
	db *sql.DB
}

func newDbConn(connURL string) (*dbService, error) {
	var dbs dbService
	var err error
	dbs.db, err = sql.Open("libsql", connURL)
	if err != nil {
		return nil, err
	}
	return &dbs, nil
}

func (dbs dbService) closeConn() error {
	err := dbs.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (dbs dbService) insertUser(username, password string) (string, error) {
	userId := generateRandomString(40)
	_, err := dbs.db.Exec("INSERT INTO users(userId, username, passwordHash) VALUES ((?), (?), (?));", userId, username, password)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (dbs dbService) checkUsernameExists(username string) bool {
	row := dbs.db.QueryRow("SELECT userId FROM users WHERE username=(?)", username)
	s := ""
	row.Scan(&s)
	return s != ""
}

func (dbs dbService) getUserId(username, password string) (string, error) {
	rows, err := dbs.db.Query("SELECT userId FROM users WHERE username=(?) AND passwordHash=(?)", username, password)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var s string
	for rows.Next() {
		rows.Scan(&s)
	}
	return s, nil
}

func (dbs dbService) setSessionId(userId, sessionId string) error {
	_, err := dbs.db.Exec("UPDATE users SET sessionId=(?) WHERE userId=(?)", sessionId, userId)
	if err != nil {
		return err
	}
	return nil
}
