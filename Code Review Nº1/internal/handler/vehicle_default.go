package handler

import (
	"Code_Review_N_1/internal"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// VehicleJSON is an struct that represents a vehicle in json format.
type VehicleJSON struct {
	ID           int     `json:"id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	Registration string  `json:"registration"`
	Year         int     `json:"year"`
	Color        string  `json:"color"`
	MaxSpeed     int     `json:"max_speed"`
	FuelType     string  `json:"fuel_type"`
	Transmission string  `json:"transmission"`
	Passengers   int     `json:"passengers"`
	Height       float64 `json:"height"`
	Width        float64 `json:"width"`
	Weight       float64 `json:"weight"`
}

// NewVehicleDefault returns a new instance of a vehicle handler.
func NewVehicleDefault(sv internal.ServiceVehicle) *VehicleDefault {
	return &VehicleDefault{sv: sv}
}

// VehicleDefault is an struct that contains handlers for vehicle.
type VehicleDefault struct {
	sv internal.ServiceVehicle
}

func convertVehiclesToJSON(vehicles []internal.Vehicle) []VehicleJSON {
	data := make([]VehicleJSON, len(vehicles))
	for i, vehicle := range vehicles {
		data[i] = VehicleJSON{
			ID:           vehicle.ID,
			Brand:        vehicle.Attributes.Brand,
			Model:        vehicle.Attributes.Model,
			Registration: vehicle.Attributes.Registration,
			Year:         vehicle.Attributes.Year,
			Color:        vehicle.Attributes.Color,
			MaxSpeed:     vehicle.Attributes.MaxSpeed,
			FuelType:     vehicle.Attributes.FuelType,
			Transmission: vehicle.Attributes.Transmission,
			Passengers:   vehicle.Attributes.Passengers,
			Height:       vehicle.Attributes.Height,
			Width:        vehicle.Attributes.Width,
			Weight:       vehicle.Attributes.Weight,
		}
	}
	return data
}

// GetAll returns all vehicles.
func (c *VehicleDefault) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// ...
		// process
		// - get all vehicles from the service
		vehicles, err := c.sv.FindAll()
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceVehicleNotFound):
				ctx.JSON(http.StatusNotFound, map[string]any{"message": "vehicles not found"})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]any{"message": "internal server error"})
			}
			return
		}

		// response
		// - serialize vehicles
		data := convertVehiclesToJSON(vehicles)
		ctx.JSON(http.StatusOK, map[string]any{"message": "success to find vehicles", "data": data})
	}
}

func (c *VehicleDefault) AddVehicle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newVehicle internal.Vehicle

		if err := ctx.BindJSON(&newVehicle); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if err := c.sv.ValidateVehicleFields(newVehicle); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := c.sv.ValidateDateFormat(newVehicle.Attributes.Registration); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := c.sv.ValidateUniqueRegistration(newVehicle.Attributes.Registration); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := c.sv.AddVehicle(newVehicle); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add vehicle"})
			return
		}

		ctx.IndentedJSON(http.StatusCreated, newVehicle)
	}
}

func (c *VehicleDefault) FindByColorAndYear() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		color := ctx.Param("color")
		year, err := strconv.Atoi(ctx.Param("year"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}

		vehicles, err := c.sv.FindByColorAndYear(color, year)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceVehicleNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"message": "No vehicles found with the specified criteria"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			}
			return
		}

		data := convertVehiclesToJSON(vehicles)
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "data": data})
	}
}

func (c *VehicleDefault) FindByBrandAndYearRange() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		brand := ctx.Param("brand")
		startYear, err := strconv.Atoi(ctx.Param("start_year"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_year"})
			return
		}
		endYear, err := strconv.Atoi(ctx.Param("end_year"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_year"})
			return
		}
		vehicles, err := c.sv.FindByBrandAndYearRange(brand, startYear, endYear)
		if err != nil {
			if errors.Is(err, internal.ErrServiceVehicleNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "vehicles not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}

		data := convertVehiclesToJSON(vehicles)
		ctx.JSON(http.StatusOK, map[string]interface{}{"message": "success to find vehicles", "data": data})
	}
}

func (c *VehicleDefault) GetAverageSpeedByBrand() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		brand := ctx.Param("brand")

		averageSpeed, err := c.sv.GetAverageSpeedByBrand(brand)
		if err != nil {
			if errors.Is(err, internal.ErrRepositoryVehicleNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "vehicles not found for the brand"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"brand": brand, "average_speed": averageSpeed})
	}
}

func (c *VehicleDefault) AddMultipleVehicles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newVehicles []internal.Vehicle

		if err := ctx.BindJSON(&newVehicles); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if err := c.sv.AddMultipleVehicles(newVehicles); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Vehicles created successfully"})
	}
}

func (c *VehicleDefault) UpdateMaxSpeed() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		newMaxSpeed := ctx.PostForm("new_max_speed")

		vehicleID, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
			return
		}

		maxSpeed, err := strconv.Atoi(newMaxSpeed)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maximum speed"})
			return
		}
		if err := c.sv.UpdateMaxSpeed(vehicleID, maxSpeed); err != nil {
			if errors.Is(err, internal.ErrRepositoryVehicleNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maximum speed"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Maximum speed updated successfully"})
	}
}

func (c *VehicleDefault) GetByFuelType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fuelType := ctx.Param("type")

		vehicles, err := c.sv.FindByFuelType(fuelType)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceVehicleNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"message": "No vehicles found with the specified fuel type"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			}
			return
		}

		data := convertVehiclesToJSON(vehicles)
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "data": data})
	}
}

func (c *VehicleDefault) DeleteVehicle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		vehicleID, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
			return
		}

		if err := c.sv.DeleteVehicleByID(vehicleID); err != nil {
			switch {
			case errors.Is(err, internal.ErrRepositoryVehicleNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle"})
			}
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}
