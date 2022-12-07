package Database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./quotes.db")
	if err != nil {
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS authors (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) UNIQUE,id INTEGER PRIMARY KEY AUTOINCREMENT ,name CHAR(60) NOT NULL check(LENGTH(name) BETWEEN 3 AND 60))")
	if err != nil {
		return nil, err
	}
	statement.Exec()
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS quotes (uuid TEXT check(uuid IS NULL OR LENGTH(uuid) > 36 OR LENGTH(uuid) < 36) ,id INTEGER PRIMARY KEY AUTOINCREMENT , quote TEXT check(quote IS NULL OR LENGTH(quote) > 300) UNIQUE, author_uuid check(author_uuid IS NULL OR LENGTH(author_uuid) > 36 OR LENGTH(author_uuid) < 36), FOREIGN KEY (author_uuid) REFERENCES authors(uuid))")
	if err != nil {
		return nil, err
	}
	statement.Exec()

	return database, nil
}

func CleanUp() error {
	err := os.Remove("./quotes.db")
	return err
}
