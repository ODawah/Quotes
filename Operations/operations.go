package Operations

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/awesomeQuotes/Schemas"
)

func InsertAuthor(db *sql.DB, author Schemas.Author) (*Schemas.Author, error) {
	err := ValidateName(author.Name)
	if err != nil {
		return nil, err
	}
	author.Name = strings.TrimSpace(strings.ToLower(author.Name))
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

func SearchAuthor(db *sql.DB, name string) (*Schemas.Author, error) {
	var author Schemas.Author
	err := ValidateName(name)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(strings.ToLower(name))
	statement, err := db.Prepare("SELECT * FROM authors WHERE name LIKE ? ")
	if err != nil {
		return nil, err
	}
	err = statement.QueryRow(name).Scan(&author.UUID, &author.ID, &author.Name)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func SearchAuthorByUUID(db *sql.DB, uuid string) (*Schemas.Author, error) {
	var author Schemas.Author
	valid := ValidateUUID(uuid)
	if !valid {
		return nil, errors.New("not valid uuid entered")
	}
	statement, err := db.Prepare("SELECT * FROM authors WHERE uuid = ? ")
	if err != nil {
		return nil, err
	}
	err = statement.QueryRow(uuid).Scan(&author.UUID, &author.ID, &author.Name)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func InsertQuote(db *sql.DB, quote Schemas.Quote) (*Schemas.Quote, error) {
	err := ValidateQuote(quote.Text)
	if err != nil {
		return nil, err
	}
	var author *Schemas.Author
	author, err = SearchAuthor(db, quote.Author.Name)
	if err != nil || author == nil {
		author, err = InsertAuthor(db, quote.Author)
		if err != nil {
			return nil, err
		}
	}
	quote.Author = *author
	quote.Text = strings.TrimSpace(strings.ToLower(quote.Text))
	uuid := UuidGenerator()
	statement, err := db.Exec("INSERT INTO quotes(uuid, quote, author_uuid) VALUES (?, ?, ?)", uuid, quote.Text, quote.Author.UUID)
	if err != nil {
		return nil, err
	}
	id, err := statement.LastInsertId()
	if err != nil || id == 0 {
		return nil, err
	}
	quote.UUID = uuid
	quote.ID = int(id)
	return &quote, nil
}

func SearchQuote(db *sql.DB, text string) (*Schemas.Quote, error) {
	var quote Schemas.Quote
	err := ValidateQuote(text)
	if err != nil {
		return nil, err
	}
	text = strings.TrimSpace(strings.ToLower(text))
	statement, err := db.Prepare("SELECT * FROM quotes WHERE quote LIKE ? ")
	if err != nil {
		return nil, err
	}
	err = statement.QueryRow(text).Scan(&quote.UUID, &quote.ID, &quote.Text, &quote.Author.UUID)
	if err != nil {
		return nil, err
	}
	author, err := SearchAuthorByUUID(db, quote.Author.UUID)
	if err != nil {
		return nil, err
	}
	quote.Author = *author
	return &quote, nil
}

func AuthorQuotes(db *sql.DB, name string) (*Schemas.QuoteList, error) {
	var result Schemas.QuoteList
	err := ValidateName(name)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(strings.ToLower(name))
	author, err := SearchAuthor(db, name)
	if err != nil || author.ID == 0 {
		return nil, err
	}
	result.Author = *author
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM quotes WHERE author_uuid = \"%s\" ", author.UUID))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var q Schemas.Quote
		rows.Scan(&q.UUID, &q.ID, &q.Text, &q.Author.UUID)
		q.Author.ID = author.ID
		result.Quotes = append(result.Quotes, q)
	}

	return &result, nil
}
