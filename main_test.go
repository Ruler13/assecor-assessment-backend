package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPersons(t *testing.T) {
	repo := NewCSVRepository()
	server := NewServer(repo)

	req, err := http.NewRequest("GET", "/persons", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.getPersons(w, r)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"Hans","lastname":"Müller","zipcode":"67742","city":"Lauterecken","color":"blau"}`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestGetPersonByID(t *testing.T) {
	repo := NewCSVRepository()
	server := NewServer(repo)

	req, err := http.NewRequest("GET", "/persons/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.getPersonByID(w, r)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":1,"name":"Hans","lastname":"Müller","zipcode":"67742","city":"Lauterecken","color":"blau"}`
	if rr.Body.String() != expected+"\n" {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetPersonsByColor(t *testing.T) {
	repo := NewCSVRepository()
	server := NewServer(repo)

	req, err := http.NewRequest("GET", "/persons/color/blau", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.getPersonsByColor(w, r)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"color":"blau"`) {
		t.Errorf("Handler did not filter by color: got %v", body)
	}
}

func TestAddPerson(t *testing.T) {
	repo := NewCSVRepository()
	server := NewServer(repo)

	personJSON := `{"name":"Test","lastname":"User","zipcode":"12345","city":"TestCity","color":"blau"}`
	req, err := http.NewRequest("POST", "/persons", bytes.NewBufferString(personJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.addPerson(w, r)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
