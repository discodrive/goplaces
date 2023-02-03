package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	r.HandleFunc("/place/{id}", makeHTTPHandlerFunc(s.handleGetPlaceByID))

	log.Printf("Server is running on port %v", s.listenerAddr)

	http.ListenAndServe(s.listenerAddr, r)
}

func (s *APIServer) handlePlace(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetPlaces(w, r)
	case "POST":
		return s.handleCreatePlace(w, r)
	case "DELETE":
		return s.handleDeletePlace(w, r)
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
	id := mux.Vars(r)["id"]
	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, &Place{})
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
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// APIFunc is the function signature for our functions
type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

// makeHTTPHandlerFunc is a decorator for our API functions to turn them into http handlers
func makeHTTPHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
