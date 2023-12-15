package internal

import (
	"errors"
)

var (
	// ErrServiceVehicleNotFound is returned when no vehicle is found.
	ErrServiceVehicleNotFound = errors.New("service: vehicle not found")
)

// ServiceVehicle is the interface that wraps the basic methods for a vehicle service.
// - conections with external apis
// - business logic
type ServiceVehicle interface {
	// FindAll returns all vehicles
	FindAll() (v []Vehicle, err error)
	AddVehicle(newVehicle Vehicle) error
	ValidateVehicleFields(vehicle Vehicle) error
	ValidateDateFormat(date string) error
	ValidateUniqueRegistration(registration string) error
	FindByColorAndYear(color string, year int) ([]Vehicle, error)
	FindByBrandAndYearRange(brand string, startYear, endYear int) (v []Vehicle, err error)
	GetAverageSpeedByBrand(brand string) (float64, error)
	AddMultipleVehicles(newVehicles []Vehicle) error
	UpdateMaxSpeed(id int, newMaxSpeed int) error
	FindByFuelType(fuelType string) ([]Vehicle, error)
	DeleteVehicleByID(id int) error
}
