package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Person struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Zipcode  string `json:"zipcode"`
	City     string `json:"city"`
	Color    string `json:"color"`
}

type Server struct {
	repo PersonRepository
}

func NewServer(repo PersonRepository) *Server {
	return &Server{repo: repo}
}

func (s *Server) getPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	persons, err := s.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(persons)
}

func (s *Server) getPersonByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/persons/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	person, err := s.repo.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (s *Server) getPersonsByColor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	color := strings.TrimPrefix(r.URL.Path, "/persons/color/")
	persons, err := s.repo.GetByColor(color)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(persons)
}

func (s *Server) addPerson(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if person.Name == "" || person.Lastname == "" || person.Zipcode == "" || person.City == "" || person.Color == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	person, err = s.repo.Add(person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func main() {
	repo := NewCSVRepository()
	server := NewServer(repo)

	http.HandleFunc("/persons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			server.addPerson(w, r)
		} else {
			server.getPersons(w, r)
		}
	})
	http.HandleFunc("/persons/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/persons/color/") {
			server.getPersonsByColor(w, r)
		} else {
			server.getPersonByID(w, r)
		}
	})

	fmt.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
