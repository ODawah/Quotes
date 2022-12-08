package Operations

import (
	"database/sql"
	"errors"
	"fmt"
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
	if name == "" {
		return nil, errors.New("no name entered")
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
	if uuid == "" {
		return nil, errors.New("no uuid entered")
	}
	statement, err := db.Prepare("SELECT * FROM authors WHERE uuid LIKE ? ")
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
	if quote.Text == "" {
		return nil, errors.New("no quote inserted")
	} else if len(quote.Text) > 300 {
		return nil, errors.New("long quote")
	} else if len(quote.Text) < 10 {
		return nil, errors.New("short quote")
	}
	var author *Schemas.Author
	author, err := SearchAuthor(db, quote.Author.Name)
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
	if text == "" {
		return nil, errors.New("no quote entered")
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
	if name == "" {
		return nil, errors.New("no author name entered")
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
		result.Quotes = append(result.Quotes, q)
	}

	return &result, nil
}
