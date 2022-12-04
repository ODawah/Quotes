package Database

import (
	"os"
	"testing"
)

func TestCleanUp(t *testing.T) {
	_, err := Connect()
	if err != nil {
		t.Fatal(err)
	}
	CleanUp()
	if _, err := os.Stat("./quotes.db"); err == nil {
		t.Fatal("Database file not deleted")
	}
}

func TestConnect(t *testing.T) {
	database, err := Connect()
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.Query("SELECT * FROM authors")
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.Query("SELECT * FROM quotes")
	if err != nil {
		t.Fatal(err)
	}
	CleanUp()
}
