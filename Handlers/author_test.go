package Handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/awesomeQuotes/Database"
	"github.com/awesomeQuotes/Operations"
	"github.com/awesomeQuotes/Schemas"
)

func TestSearchAuthorHandler(t *testing.T) {
	db, _ := Database.Connect()
	Operations.InsertAuthor(db, Schemas.Author{Name: "omar"})
	defer Database.CleanUp()

	type test struct {
		name     string
		url      string
		expected *Schemas.Author
		err      bool
	}

	tests := []test{
		{"normal name", "http://localhost:8080/find_author/omar", &Schemas.Author{Name: "omar", ID: 1}, false},
		{"capital name", "http://localhost:8080/find_author/OMAR", &Schemas.Author{Name: "omar", ID: 1}, false},
		{"name not in database", "http://localhost:8080/find_author/george", nil, true},
		{"no search name", "http://localhost:8080/find_author/", nil, true},
	}

	for _, tc := range tests {
		var got *Schemas.Author
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
			if got.ID != tc.expected.ID {
				Database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.ID, tc.expected.ID)
			}
		}

	}

}

func TestCreateAuthor(t *testing.T) {
	defer Database.CleanUp()
	type test struct {
		name     string
		input    *Schemas.Author
		expected *Schemas.Author
		err      bool
	}
	tests := []test{
		{"normal name", &Schemas.Author{Name: "william"}, &Schemas.Author{Name: "william", ID: 1}, false},
		{"capital name", &Schemas.Author{Name: "JENNIFER"}, &Schemas.Author{Name: "jennifer", ID: 2}, false},
		{"empty body", &Schemas.Author{}, nil, true}, {"no body", nil, nil, true},
	}
	for _, tc := range tests {
		Database.Connect()
		var got *Schemas.Author
		var buffer bytes.Buffer
		err := json.NewEncoder(&buffer).Encode(tc.input)
		if err != nil {
			t.Log(tc.name)
			t.Log("error encoding json")
			t.Fatal(err)
		}
		resp, err := http.Post("http://localhost:8080/create_author", "application/json", &buffer)
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
			if got.Name != tc.expected.Name {
				t.Logf("test name: %s", tc.name)
				t.Fatalf("got: %s  expected: %s", got.Name, tc.expected.Name)
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
