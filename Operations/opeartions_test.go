package Operations

import (
	"testing"

	"github.com/awesomeQuotes/Database"
	"github.com/awesomeQuotes/Schemas"
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
	db, _ := Database.Connect()
	type test struct {
		name     string
		input    Schemas.Author
		expected *Schemas.Author
		err      bool
	}
	tests := []test{
		{"normal name", Schemas.Author{Name: "omar"}, &Schemas.Author{Name: "omar", ID: 1}, false},
		{"capital name", Schemas.Author{Name: "ADHAM"}, &Schemas.Author{Name: "adham", ID: 2}, false},
		{"name less than 3 chars", Schemas.Author{Name: "om"}, nil, true},
		{"no name", Schemas.Author{Name: ""}, nil, true},
		{"input larger than constraint (70)", Schemas.Author{Name: "john ben karim samir george johnny omar ahmed mahmoud masouds"}, nil, true},
	}

	for _, tc := range tests {
		got, err := InsertAuthor(db, tc.input)
		if (err != nil) != tc.err {
			Database.CleanUp()
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Name != tc.expected.Name {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Name, tc.expected.Name)
			}
			if len(got.UUID) != 36 {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
		}
	}
	Database.CleanUp()
}

func TestSearchAuthor(t *testing.T) {
	db, _ := Database.Connect()

	authors := []Schemas.Author{
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
		expected *Schemas.Author
		err      bool
	}
	tests := []test{
		{"normal name", "omar", &Schemas.Author{Name: "omar", ID: 1}, false},
		{"capital name", "ADHAM", &Schemas.Author{Name: "adham", ID: 2}, false},
		{"name not in database", "george", nil, true},
		{"no search name", "", nil, true},
	}

	for _, tc := range tests {
		got, err := SearchAuthor(db, tc.search)
		if (err != nil) != tc.err {
			Database.CleanUp()
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Name != tc.expected.Name {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Name, tc.expected.Name)
			}
			if len(got.UUID) != 36 {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
		}
	}
	Database.CleanUp()
}

func TestInsertQuote(t *testing.T) {
	// to make the test easy to read
	longtext := "I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"
	db, _ := Database.Connect()

	// insert authors for the tests
	authors := []Schemas.Author{
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
		input    Schemas.Quote
		expected *Schemas.Quote
		err      bool
	}
	tests := []test{
		{"author in database", Schemas.Quote{Text: "heaven is for real", Author: Schemas.Author{Name: "omar"}}, &Schemas.Quote{Text: "heaven is for real", Author: Schemas.Author{Name: "omar", ID: 1}}, false},
		{"author not in database", Schemas.Quote{Text: "Keep Dreaming", Author: Schemas.Author{Name: "maged"}}, &Schemas.Quote{Text: "keep dreaming", Author: Schemas.Author{Name: "maged", ID: 3}}, false},
		{"no quote", Schemas.Quote{Text: "", Author: Schemas.Author{Name: "maged"}}, nil, true},
		{"long quote", Schemas.Quote{Text: longtext, Author: Schemas.Author{Name: "omar"}}, nil, true},
		{"short quote", Schemas.Quote{Text: "short", Author: Schemas.Author{Name: "omar"}}, nil, true},
	}

	for _, tc := range tests {
		got, err := InsertQuote(db, tc.input)
		if (err != nil) != tc.err {
			Database.CleanUp()
			t.Logf("test name: %s", tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Text != tc.expected.Text {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Text, tc.expected.Text)
			}
			if len(got.UUID) != 36 {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if len(got.Author.UUID) != 36 {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.Author.UUID)
			}
			if got.Author.Name != tc.expected.Author.Name {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Author.Name, tc.expected.Author.Name)

			}
			if got.Author.ID != tc.expected.Author.ID {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.Author.ID, tc.expected.Author.ID)

			}
		}
	}
	Database.CleanUp()
}
