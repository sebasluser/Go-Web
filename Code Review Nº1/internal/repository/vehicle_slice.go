package repository

import "Code_Review_N_1/internal"

// NewVehicleSlice returns a new instance of a vehicle repository in an slice.
func NewVehicleSlice(db []internal.Vehicle, lastId int) *VehicleSlice {
	return &VehicleSlice{
		db:     db,
		lastId: lastId,
	}
}

// VehicleSlice is an struct that represents a vehicle repository in an slice.
type VehicleSlice struct {
	// db is the database of vehicles.
	db []internal.Vehicle
	// lastId is the last id of the database.
	lastId int
}

// FindAll returns all vehicles
func (s *VehicleSlice) FindAll() (v []internal.Vehicle, err error) {
	// check if the database is empty
	if len(s.db) == 0 {
		err = internal.ErrRepositoryVehicleNotFound
		return
	}

	// make a copy of the database
	v = make([]internal.Vehicle, len(s.db))
	copy(v, s.db)
	return
}

func (s *VehicleSlice) AddVehicle(newVehicles internal.Vehicle) error {
	s.db = append(s.db, newVehicles)
	return nil
}

func (s *VehicleSlice) AddMultipleVehicles(newVehicles []internal.Vehicle) error {
	for _, vehicle := range newVehicles {
		if err := s.AddVehicle(vehicle); err != nil {
			return err
		}
	}
	return nil
}

func (s *VehicleSlice) UpdateMaxSpeed(id int, newMaxSpeed int) error {
	for i, vehicle := range s.db {
		if vehicle.ID == id {
			s.db[i].Attributes.MaxSpeed = newMaxSpeed
			return nil
		}
	}
	return internal.ErrRepositoryVehicleNotFound
}
func (s *VehicleSlice) DeleteByID(id int) error {
	index := -1
	for i, vehicle := range s.db {
		if vehicle.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return internal.ErrRepositoryVehicleNotFound
	}
	s.db = append(s.db[:index], s.db[index+1:]...)

	return nil
}
