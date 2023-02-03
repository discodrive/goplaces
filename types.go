package main

import (
	"time"
)

// Best practice to not use Place for requests because we will not be providing an ID
type CreatePlaceRequest struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

// Place represents data about a location
type Place struct {
	ID        int       `json:"id"`
	Location  string    `json:"location"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewPlace(location string, name string) *Place {
	return &Place{
		Location:  location,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
}
