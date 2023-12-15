package service

import (
	"Code_Review_N_1/internal"
	"errors"
	"fmt"
	"time"
)

// NewDefault returns a new instance of a vehicle service.
func NewDefault(rp internal.RepositoryVehicle) *Default {
	return &Default{rp: rp}
}

// Default is an struct that represents a vehicle service.
type Default struct {
	rp internal.RepositoryVehicle
}

// FindAll returns all vehicles.
func (s *Default) FindAll() (v []internal.Vehicle, err error) {
	// get all vehicles from the repository
	v, err = s.rp.FindAll()
	if err != nil {
		if errors.Is(err, internal.ErrRepositoryVehicleNotFound) {
			err = fmt.Errorf("%w. %v", internal.ErrServiceVehicleNotFound, err)
			return
		}
		return
	}
	return
}

func (s *Default) AddVehicle(newVehicle internal.Vehicle) error {
	return s.rp.AddVehicle(newVehicle)
}

func (s *Default) ValidateVehicleFields(vehicle internal.Vehicle) error {
	if vehicle.Attributes.Brand == "" || vehicle.Attributes.Model == "" || vehicle.Attributes.Registration == "" ||
		vehicle.Attributes.Year == 0 || vehicle.Attributes.Color == "" || vehicle.Attributes.MaxSpeed == 0 ||
		vehicle.Attributes.FuelType == "" || vehicle.Attributes.Transmission == "" || vehicle.Attributes.Passengers == 0 ||
		vehicle.Attributes.Height == 0 || vehicle.Attributes.Width == 0 || vehicle.Attributes.Weight == 0 {
		return fmt.Errorf("all fields must be provided")
	}
	return nil
}

func (s *Default) ValidateDateFormat(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}
	return nil
}

func (s *Default) ValidateUniqueRegistration(registration string) error {
	// Implement the logic to check uniqueness in your repository
	// Example: Check if the registration already exists in the repository
	vehicles, err := s.rp.FindAll()
	if err != nil {
		return err
	}
	for _, existingVehicle := range vehicles {
		if existingVehicle.Attributes.Registration == registration {
			return fmt.Errorf("registration must be unique")
		}
	}
	return nil
}

func (s *Default) FindByColorAndYear(color string, year int) ([]internal.Vehicle, error) {
	vehicles, err := s.rp.FindAll()
	if err != nil {
		return nil, err
	}

	var matchingVehicles []internal.Vehicle
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Color == color && vehicle.Attributes.Year == year {
			matchingVehicles = append(matchingVehicles, vehicle)
		}
	}

	if len(matchingVehicles) == 0 {
		return nil, internal.ErrServiceVehicleNotFound
	}

	return matchingVehicles, nil
}

func (s *Default) FindByBrandAndYearRange(brand string, startYear, endYear int) (v []internal.Vehicle, err error) {
	vehicles, err := s.rp.FindAll()
	if err != nil {
		if errors.Is(err, internal.ErrRepositoryVehicleNotFound) {
			err = fmt.Errorf("%w. %v", internal.ErrServiceVehicleNotFound, err)
			return
		}
		return
	}
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Brand == brand && vehicle.Attributes.Year >= startYear && vehicle.Attributes.Year <= endYear {
			v = append(v, vehicle)
		}
	}
	if len(v) == 0 {
		err = internal.ErrServiceVehicleNotFound
		return
	}
	return v, nil
}

func (s *Default) GetAverageSpeedByBrand(brand string) (float64, error) {
	vehicles, err := s.rp.FindAll()
	if err != nil {
		return 0, err
	}

	var totalSpeed, count int
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Brand == brand {
			totalSpeed += vehicle.Attributes.MaxSpeed
			count++
		}
	}

	if count == 0 {
		return 0, internal.ErrRepositoryVehicleNotFound
	}

	averageSpeed := float64(totalSpeed) / float64(count)
	return averageSpeed, nil
}

func (s *Default) AddMultipleVehicles(newVehicles []internal.Vehicle) error {
	for _, vehicle := range newVehicles {
		if err := s.rp.AddVehicle(vehicle); err != nil {
			return err
		}
	}
	return nil
}

func (s *Default) UpdateMaxSpeed(id int, newMaxSpeed int) error {
	return s.rp.UpdateMaxSpeed(id, newMaxSpeed)
}

func (s *Default) FindByFuelType(fuelType string) ([]internal.Vehicle, error) {
	vehicles, err := s.rp.FindAll()
	if err != nil {
		return nil, err
	}

	var matchingVehicles []internal.Vehicle
	for _, vehicle := range vehicles {
		if vehicle.Attributes.FuelType == fuelType {
			matchingVehicles = append(matchingVehicles, vehicle)
		}
	}

	if len(matchingVehicles) == 0 {
		return nil, internal.ErrServiceVehicleNotFound
	}

	return matchingVehicles, nil
}

func (s *Default) DeleteVehicleByID(id int) error {
	if err := s.rp.DeleteByID(id); err != nil {
		return err
	}
	return nil
}
