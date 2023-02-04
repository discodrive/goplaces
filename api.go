package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenerAddr string
	store        Storage
}

// NewAPIServer returns a pointer to our API server
func NewAPIServer(listenerAddr string, store Storage) *APIServer {
	return &APIServer{
		listenerAddr: listenerAddr,
		store:        store,
	}
}

func (s *APIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/places", makeHTTPHandlerFunc(s.handlePlace))
	r.HandleFunc("/places/{id}", makeHTTPHandlerFunc(s.handleGetPlaceByID))

	log.Printf("Server is running on port %v", s.listenerAddr)

	http.ListenAndServe(s.listenerAddr, r)
}

func (s *APIServer) handlePlace(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetPlaces(w, r)
	case "POST":
		return s.handleCreatePlace(w, r)
	}
	return fmt.Errorf("Unsupported method %s", r.Method)
}

func (s *APIServer) handleGetPlaces(w http.ResponseWriter, r *http.Request) error {
	places, err := s.store.GetPlaces()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, places)
}

func (s *APIServer) handleGetPlaceByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		place, err := s.store.GetPlaceByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, place)
	}

	if r.Method == "DELETE" {
		return s.handleDeletePlace(w, r)
	}

	return fmt.Errorf("Method %s not allowed", r.Method)
}

func (s *APIServer) handleCreatePlace(w http.ResponseWriter, r *http.Request) error {
	createPlaceReq := new(CreatePlaceRequest)
	if err := json.NewDecoder(r.Body).Decode(createPlaceReq); err != nil {
		return err
	}

	place := NewPlace(createPlaceReq.Location, createPlaceReq.Name)
	if err := s.store.CreatePlace(place); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, place)
}

func (s *APIServer) handleDeletePlace(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeletePlace(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// APIFunc is the function signature for our functions
type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

// makeHTTPHandlerFunc is a decorator for our API functions to turn them into http handlers
func makeHTTPHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	// Fetch the ID from the request param as a string and convert it to int
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Return our own friendly error rather than the JSON Atoi error which isn't as helpful
		return id, fmt.Errorf("ID %s is not valid.", idStr)
	}
	return id, nil
}
