package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/awesomeQuotes/database"
	"github.com/awesomeQuotes/operations"
	"github.com/awesomeQuotes/schemas"
)

func TestSearchQuote(t *testing.T) {
	db, _ := database.Connect()
	operations.InsertQuote(db, schemas.Quote{
		Text:   "Keep it simple, stupid",
		Author: schemas.Author{Name: "omar"},
	})
	defer database.CleanUp()
	type test struct {
		name     string
		input    *schemas.Quote
		expected *schemas.Quote
		err      bool
	}
	tests := []test{
		{"normal quote", &schemas.Quote{Text: "keep it simple, stupid"}, &schemas.Quote{Text: "keep it simple, stupid", ID: 1, Author: schemas.Author{Name: "omar", ID: 1}}, false},
		{"quote not in database", &schemas.Quote{Text: "Keep dreaming"}, nil, false},
		{"empty body", nil, nil, true}, {"no body", nil, nil, true},
	}

	for _, tc := range tests {
		var got *schemas.Quote
		client := &http.Client{}
		b, err := json.Marshal(tc.input)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("GET", "http://localhost:8080/find_quote", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		req.Close = true
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(body, &got); err != nil && tc.expected != nil {
			t.Log(tc.name)
			t.Fatal(err)
		}
		if (got != nil) && (tc.expected != nil) {
			if got.Text != tc.expected.Text {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %s  expected: %s", got.Author.Name, tc.expected.Author.Name)
			}
			if got.Author.Name != tc.expected.Author.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %s  expected: %s", got.Author.Name, tc.expected.Author.Name)
			}
			if len(got.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if len(got.Author.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if got.ID != tc.expected.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %d  expected: %d", got.ID, tc.expected.ID)
			}
			if got.Author.ID != tc.expected.Author.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %d  expected: %d", got.ID, tc.expected.ID)
			}
		}

	}
}

func TestCreateQuote(t *testing.T) {
	defer database.CleanUp()
	type test struct {
		name     string
		input    *schemas.Quote
		expected *schemas.Quote
		err      bool
	}
	tests := []test{
		{"normal quote", &schemas.Quote{Text: "keep it real", Author: schemas.Author{Name: "omar", ID: 1}}, &schemas.Quote{Text: "keep it real", ID: 1, Author: schemas.Author{Name: "omar", ID: 1}}, false},
		{"quote duplicate", &schemas.Quote{Text: "keep it simple, stupid", Author: schemas.Author{Name: "omar", ID: 1}}, nil, false},
		{"no quote", &schemas.Quote{Author: schemas.Author{Name: "omar", ID: 1}}, nil, false},
		{"no quote author", &schemas.Quote{Text: "Keep dreaming"}, nil, false},
		{"empty body", nil, nil, true}, {"no body", nil, nil, true},
	}
	for _, tc := range tests {
		var got *schemas.Quote
		var buffer bytes.Buffer
		err := json.NewEncoder(&buffer).Encode(tc.input)
		if err != nil {
			t.Log(tc.name)
			t.Log("error encoding json")
			t.Fatal(err)
		}
		resp, err := http.Post("http://localhost:8080/create_quote", "application/json", &buffer)
		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(body, &got); err != nil && tc.expected != nil {
			t.Log(tc.name)
			t.Log("here")
			t.Fatal(err)
		}
		t.Log(got)
		if (got != nil) && (tc.expected != nil) {
			if got.Text != tc.expected.Text {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %s  expected: %s", got.Text, tc.expected.Text)

			}
			if got.Author.Name != tc.expected.Author.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %s  expected: %s", got.Author.Name, tc.expected.Author.Name)
			}
			if len(got.Author.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.Author.UUID)
			}
			if len(got.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if got.ID != tc.expected.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %d  expected: %d", got.ID, tc.expected.ID)
			}
		}
	}
}

func TestAuthorQuotes(t *testing.T) {
	db, _ := database.Connect()
	operations.InsertQuote(db, schemas.Quote{
		Text:   "Keep it simple, stupid",
		Author: schemas.Author{Name: "omar"},
	})
	operations.InsertQuote(db, schemas.Quote{
		Text:   "Go with the flow",
		Author: schemas.Author{Name: "omar"},
	})

	defer database.CleanUp()
	type test struct {
		name     string
		url      string
		expected *schemas.QuoteList
		err      bool
	}

	tests := []test{
		{"normal name", "http://localhost:8080/find_Author_quotes/omar", &schemas.QuoteList{Author: schemas.Author{Name: "omar", ID: 1}, Quotes: []schemas.Quote{
			{Text: "keep it simple, stupid", Author: schemas.Author{Name: "omar", ID: 1}},
			{Text: "go with the flow", Author: schemas.Author{Name: "omar", ID: 1}},
		}}, false},
		{"name not in database", "http://localhost:8080/find_Author_quotes/george", nil, true},
		{"no search name", "http://localhost:8080/find_Author_quotes/", nil, true},
	}

	for _, tc := range tests {
		var got *schemas.QuoteList
		resp, err := http.Get(tc.url)
		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(body, &got); err != nil && tc.expected != nil {
			t.Log(tc.name)
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
			if got.Author.Name != tc.expected.Author.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %s  expected: %s", got.Author.Name, tc.expected.Author.Name)
			}
			if len(got.Author.UUID) != 36 {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.Author.UUID)
			}
			if got.Author.ID != tc.expected.Author.ID {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %d  expected: %d", got.Author.ID, tc.expected.Author.ID)
			}

		}

	}

}
