package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func setup() {
	godotenv.Load()
}

func TestReposHelper(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/repos", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ReposHelper(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that we don't have empty map
	var data []interface{}
	if json.Unmarshal([]byte(rr.Body.String()), &data); len(data) == 0 {
		t.Errorf("returned empty map of data for the specified repositories")
	}
}

func TestFetchResource(t *testing.T) {

}
