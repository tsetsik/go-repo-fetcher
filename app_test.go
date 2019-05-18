package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReposHelper(t *testing.T) {
	req, err := http.NewRequest("GET", "https://api.github.com/users/vmware/repos", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReposHelper)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that we have more than we have not empty map
	var data []interface{}
	if json.Unmarshal([]byte(rr.Body.String()), &data); len(data) == 0 {
		t.Errorf("returned empty map of data for the specified repositories")
	}
}
