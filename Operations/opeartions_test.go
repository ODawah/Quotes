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

func TestAuthorSchemas(t *testing.T) {
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
			if tc.expected != nil {
				if len(got.UUID) != 36 {
					Database.CleanUp()
					t.Logf("test name: %s", tc.name)
					t.Fatalf("worng uuid: %s", got.UUID)
				}
			}
		}
	}
	Database.CleanUp()
}
