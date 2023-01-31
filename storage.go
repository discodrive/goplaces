package main

type Storage interface {
	CreatePlace(*Place) error
	DeletePlace(int) error
	UpdatePlace(*Place) error
	GetPlaceByID(int) (*Place, error)
}
