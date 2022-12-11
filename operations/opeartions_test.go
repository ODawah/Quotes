package operations

import (
	"testing"

	"github.com/awesomeQuotes/database"
	"github.com/awesomeQuotes/schemas"
)

func TestUuidGenerator(t *testing.T) {
	uuid := UuidGenerator()
	if uuid == "" || len(uuid) != 36 {
		t.Log(len(uuid))
		t.Log(uuid)
		t.Fatal("UUID is empty")
	}
}

func TestInsertAuthor(t *testing.T) {
	defer database.CleanUp()
	db, _ := database.Connect()
	type test struct {
		name     string
		input    schemas.Author
		expected *schemas.Author
		err      bool
	}
	tests := []test{
		{"normal name", schemas.Author{Name: "omar"}, &schemas.Author{Name: "omar", ID: 1}, false},
		{"capital name", schemas.Author{Name: "ADHAM"}, &schemas.Author{Name: "adham", ID: 2}, false},
		{"name less than 3 chars", schemas.Author{Name: "om"}, nil, true},
		{"no name", schemas.Author{Name: ""}, nil, true},
		{"input larger than constraint (70)", schemas.Author{Name: "john ben karim samir george johnny omar ahmed mahmoud masouds"}, nil, true},
	}

	for _, tc := range tests {
		got, err := InsertAuthor(db, tc.input)
		if (err != nil) != tc.err {
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Name != tc.expected.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Name, tc.expected.Name)
			}
			if len(got.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if got.ID != tc.expected.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.ID, tc.expected.ID)
			}
		}
	}
}

func TestSearchAuthor(t *testing.T) {
	db, _ := database.Connect()
	defer database.CleanUp()

	authors := []schemas.Author{
		{Name: "omar"},
		{Name: "adham"},
		{Name: "maged"},
	}

	for _, tc := range authors {
		_, err := InsertAuthor(db, tc)
		if err != nil {
			t.Fatalf("error inserting authors: %s ", err)
		}
	}

	type test struct {
		name     string
		search   string
		expected *schemas.Author
		err      bool
	}
	tests := []test{
		{"normal name", "omar", &schemas.Author{Name: "omar", ID: 1}, false},
		{"capital name", "ADHAM", &schemas.Author{Name: "adham", ID: 2}, false},
		{"name not in database", "george", nil, true},
		{"no search name", "", nil, true},
	}

	for _, tc := range tests {
		got, err := SearchAuthor(db, tc.search)
		if (err != nil) != tc.err {
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Name != tc.expected.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Name, tc.expected.Name)
			}
			if len(got.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if got.ID != tc.expected.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.ID, tc.expected.ID)
			}
		}
	}
}

func TestSearchAuthorByUUID(t *testing.T) {
	db, _ := database.Connect()
	defer database.CleanUp()
	var uuidList []string
	authors := []schemas.Author{
		{Name: "omar"},
		{Name: "adham"},
		{Name: "maged"},
	}

	for _, tc := range authors {
		author, err := InsertAuthor(db, tc)
		if err != nil {
			t.Fatalf("error inserting authors: %s ", err)
		}
		uuidList = append(uuidList, author.UUID)
	}

	type test struct {
		name     string
		uuid     string
		expected *schemas.Author
		err      bool
	}
	tests := []test{
		{"author UUID in database", uuidList[0], &schemas.Author{Name: "omar", ID: 1}, false},
		{"author UUID not in database", "aasdweqd1-aseqweg3-qe120oe-owek1olsd", nil, true},
		{"no UUID", "", nil, true},
	}

	for _, tc := range tests {
		got, err := SearchAuthorByUUID(db, tc.uuid)
		if (err != nil) != tc.err {
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Name != tc.expected.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Name, tc.expected.Name)
			}
			if len(got.UUID) != 36 || got.UUID != tc.uuid {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got uuid: %s   expected uuid: %s", got.UUID, tc.uuid)
			}
			if got.ID != tc.expected.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.ID, tc.expected.ID)
			}
		}
	}
}

func TestInsertQuote(t *testing.T) {
	// to make the test easy to read
	longtext := "I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"
	db, _ := database.Connect()
	defer database.CleanUp()

	// insert authors for the tests
	authors := []schemas.Author{
		{Name: "omar"},
		{Name: "adham"},
	}
	for _, tc := range authors {
		_, err := InsertAuthor(db, tc)
		if err != nil {
			t.Fatalf("error inserting authors: %s ", err)
		}
	}

	type test struct {
		name     string
		input    schemas.Quote
		expected *schemas.Quote
		err      bool
	}
	tests := []test{
		{"author in database", schemas.Quote{Text: "heaven is for real", Author: schemas.Author{Name: "omar"}}, &schemas.Quote{Text: "heaven is for real", Author: schemas.Author{Name: "omar", ID: 1}}, false},
		{"author not in database", schemas.Quote{Text: "Keep Dreaming", Author: schemas.Author{Name: "maged"}}, &schemas.Quote{Text: "keep dreaming", Author: schemas.Author{Name: "maged", ID: 3}}, false},
		{"no quote", schemas.Quote{Text: "", Author: schemas.Author{Name: "maged"}}, nil, true},
		{"long quote", schemas.Quote{Text: longtext, Author: schemas.Author{Name: "omar"}}, nil, true},
		{"short quote", schemas.Quote{Text: "short", Author: schemas.Author{Name: "omar"}}, nil, true},
	}

	for _, tc := range tests {
		got, err := InsertQuote(db, tc.input)
		if (err != nil) != tc.err {
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Text != tc.expected.Text {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Text, tc.expected.Text)
			}
			if len(got.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if len(got.Author.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.Author.UUID)
			}
			if got.Author.Name != tc.expected.Author.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Author.Name, tc.expected.Author.Name)

			}
			if got.Author.ID != tc.expected.Author.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.Author.ID, tc.expected.Author.ID)

			}
		}
	}
}

func TestSearchQuote(t *testing.T) {
	db, _ := database.Connect()
	defer database.CleanUp()

	// insert quotes for the tests
	quotes := []schemas.Quote{
		{Text: "heaven is for real", Author: schemas.Author{Name: "omar"}},
		{Text: "keep dreaming", Author: schemas.Author{Name: "omar"}},
		{Text: "work hard and non stop", Author: schemas.Author{Name: "adham"}},
	}

	for _, tc := range quotes {
		_, err := InsertQuote(db, tc)
		if err != nil {
			t.Fatalf("error inserting quotes: %s ", err)
		}
	}

	type test struct {
		name        string
		searchQuote string
		expected    *schemas.Quote
		err         bool
	}
	tests := []test{
		{"quote in database", "heaven is for real", &schemas.Quote{Text: "heaven is for real", Author: schemas.Author{Name: "omar", ID: 1}}, false},
		{"quote not in database", "get yourself out mud", nil, true},
		{"no quote", "", nil, true},
	}
	for _, tc := range tests {
		got, err := SearchQuote(db, tc.searchQuote)
		if (err != nil) != tc.err {
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Text != tc.expected.Text {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Text, tc.expected.Text)
			}
			if len(got.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if len(got.Author.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.Author.UUID)
			}
			if got.Author.Name != tc.expected.Author.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Author.Name, tc.expected.Author.Name)

			}
			if got.Author.ID != tc.expected.Author.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.Author.ID, tc.expected.Author.ID)

			}
		}
	}
}

func TestAuthorQuotes(t *testing.T) {
	db, _ := database.Connect()
	defer database.CleanUp()

	InsertAuthor(db, schemas.Author{Name: "hany"})
	// insert quotes for the tests
	quotes := []schemas.Quote{
		{Text: "heaven is for real", Author: schemas.Author{Name: "omar"}},
		{Text: "keep dreaming", Author: schemas.Author{Name: "omar"}},
		{Text: "work hard and non stop", Author: schemas.Author{Name: "adham"}},
	}

	for _, tc := range quotes {
		_, err := InsertQuote(db, tc)
		if err != nil {
			t.Fatalf("error inserting quotes: %s ", err)
		}
	}

	type test struct {
		name         string
		searchAuthor string
		expected     *schemas.QuoteList
		err          bool
	}

	tests := []test{
		{"author in database", "omar", &schemas.QuoteList{Author: schemas.Author{Name: "omar", ID: 2}, Quotes: []schemas.Quote{
			{Text: "heaven is for real", Author: schemas.Author{Name: "omar", ID: 2}},
			{Text: "keep dreaming", Author: schemas.Author{Name: "omar", ID: 2}},
		}}, false},
		{"author with no quotes", "hany", &schemas.QuoteList{}, false},
		{"no author", "", nil, true},
	}

	for _, tc := range tests {
		got, err := AuthorQuotes(db, tc.searchAuthor)
		if (err != nil) != tc.err {
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		t.Log(got)
		if (got != nil) && (tc.expected != nil) {
			for i, quote := range got.Quotes {
				if quote.Text != tc.expected.Quotes[i].Text {
					t.Logf("test name: %s", tc.name)
					t.Fatalf("got: %s  expected: %s", quote.Text, tc.expected.Quotes[i].Text)
				}
				if len(quote.UUID) != 36 {
					t.Logf("test name: %s", tc.name)
					t.Fatalf("worng uuid: %s", quote.UUID)
				}
				if len(quote.Author.UUID) != 36 {
					t.Logf("test name: %s", tc.name)
					t.Fatalf("worng uuid: %s", quote.Author.UUID)
				}
				if got.Author.Name != tc.expected.Author.Name {
					t.Logf("test name: %s", tc.name)
					t.Fatalf("got: %s  expected: %s", quote.Author.Name, tc.expected.Author.Name)

				}
				if got.Author.ID != tc.expected.Author.ID {
					t.Logf("test name: %s", tc.name)
					t.Fatalf("got: %d  expected: %d", quote.Author.ID, tc.expected.Author.ID)

				}
			}
		}
	}

}
