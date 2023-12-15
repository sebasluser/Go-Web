package internal

import (
	"errors"
)

var (
	// ErrRepositoryVehicleNotFound is returned when a vehicle is not found.
	ErrRepositoryVehicleNotFound = errors.New("repository: vehicle not found")
)

// RepositoryVehicle is the interface that wraps the basic methods for a vehicle repository.
type RepositoryVehicle interface {
	// FindAll returns all vehicles
	FindAll() (v []Vehicle, err error)
	AddVehicle(newVehicle Vehicle) error
	AddMultipleVehicles(newVehicles []Vehicle) error
	UpdateMaxSpeed(id int, newMaxSpeed int) error
	DeleteByID(id int) error
}
