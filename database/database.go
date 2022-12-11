package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "/home/lenovo/go/src/github.com/awesomeQuotes/quotes.db")
	if err != nil {
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS authors (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) UNIQUE,id INTEGER PRIMARY KEY AUTOINCREMENT ,name CHAR(60) NOT NULL check(LENGTH(name) BETWEEN 3 AND 60))")
	if err != nil {
		return nil, err
	}
	statement.Exec()
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS quotes (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) UNIQUE,id INTEGER PRIMARY KEY AUTOINCREMENT,quote TEXT NOT NULL check(LENGTH(quote) BETWEEN 10 AND 300) UNIQUE,author_uuid NOT NULL check(LENGTH(author_uuid) = 36),FOREIGN KEY (author_uuid) REFERENCES authors(uuid))")
	if err != nil {
		return nil, err
	}
	statement.Exec()

	return database, nil
}

func CleanUp() error {
	err := os.Remove("/home/lenovo/go/src/github.com/awesomeQuotes/quotes.db")
	return err
}
