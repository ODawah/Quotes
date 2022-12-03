package Database

import (
	"database/sql"
)

func Connect() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./quotes.db")
	if err != nil {
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS  authors (uuid TEXT check(uuid IS NULL OR LENGTH(uuid) > 36 OR LENGTH(uuid) < 36) ,id INTEGER PRIMARY KEY AUTOINCREMENT , name TEXT check(uuid IS NULL OR LENGTH(name) > 60))")
	if err != nil {
		return nil, err
	}
	statement.Exec()
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS  quotes (uuid TEXT check(uuid IS NULL OR LENGTH(uuid) > 36 OR LENGTH(uuid) < 36) ,id INTEGER PRIMARY KEY AUTOINCREMENT , quote TEXT check(uuid IS NULL OR LENGTH(name) > 300), author_uuidcheck(uuid IS NULL OR LENGTH(uuid) > 36 OR LENGTH(uuid) < 36), FOREIGN KEY (author_uuid) REFERENCES authors(uuid))")
	if err != nil {
		return nil, err
	}
	statement.Exec()

	return database, nil
}

func Destruct(database *sql.DB) error {
	statement, err := database.Prepare("DROP TABLE authors")
	if err != nil {
		return err
	}
	statement.Exec()
	statement, err = database.Prepare("DROP TABLE quotes")
	if err != nil {
		return err
	}
	statement.Exec()

	return nil
}
