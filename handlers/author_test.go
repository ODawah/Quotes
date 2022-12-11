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

func TestSearchAuthorHandler(t *testing.T) {
	db, _ := database.Connect()
	operations.InsertAuthor(db, schemas.Author{Name: "omar"})
	defer database.CleanUp()

	type test struct {
		name     string
		url      string
		expected *schemas.Author
		err      bool
	}

	tests := []test{
		{"normal name", "http://localhost:8080/find_author/omar", &schemas.Author{Name: "omar", ID: 1}, false},
		{"capital name", "http://localhost:8080/find_author/OMAR", &schemas.Author{Name: "omar", ID: 1}, false},
		{"name not in database", "http://localhost:8080/find_author/george", nil, true},
		{"no search name", "http://localhost:8080/find_author/", nil, true},
	}

	for _, tc := range tests {
		var got *schemas.Author
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
				database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %s  got: %s", got.Name, tc.expected.Name)
			}
			if len(got.UUID) != 36 {
				database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("worng uuid: %s", got.UUID)
			}
			if got.ID != tc.expected.ID {
				database.CleanUp()
				t.Logf("test name: %s", tc.name)
				t.Fatalf("expected: %d  got: %d", got.ID, tc.expected.ID)
			}
		}

	}

}

func TestCreateAuthor(t *testing.T) {
	defer database.CleanUp()
	type test struct {
		name     string
		input    *schemas.Author
		expected *schemas.Author
		err      bool
	}
	tests := []test{
		{"normal name", &schemas.Author{Name: "william"}, &schemas.Author{Name: "william", ID: 1}, false},
		{"capital name", &schemas.Author{Name: "JENNIFER"}, &schemas.Author{Name: "jennifer", ID: 2}, false},
		{"empty body", &schemas.Author{}, nil, true}, {"no body", nil, nil, true},
	}
	for _, tc := range tests {
		var got *schemas.Author
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
