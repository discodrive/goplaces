package main

import (
	"github.com/google/uuid"
)

// Place represents data about a location
type Place struct {
	ID       uuid.UUID `json:"id"`
	Location string    `json:"location"`
}

func NewPlace(location string) *Place {
	return &Place{
		ID:       uuid.New(),
		Location: location,
	}
}
