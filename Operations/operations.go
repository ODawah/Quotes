package Operations

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/awesomeQuotes/Schemas"
)

func InsertAuthor(db *sql.DB, author Schemas.Author) (*Schemas.Author, error) {
	if author.Name == "" {
		return nil, errors.New("no name inserted")
	} else if len(author.Name) > 60 {
		return nil, errors.New("long name")
	} else if len(author.Name) < 3 {
		return nil, errors.New("short name")
	}
	author.Name = strings.ToLower(author.Name)
	uuid := UuidGenerator()
	statement, err := db.Exec("INSERT INTO authors(uuid, name) VALUES (?, ?)", uuid, author.Name)
	if err != nil {
		return nil, err
	}
	id, err := statement.LastInsertId()
	if err != nil || id == 0 {
		return nil, err
	}
	author.UUID = uuid
	author.ID = int(id)
	return &author, nil
}