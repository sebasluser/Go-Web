package application

import (
	"Code_Review_N_1/internal/handler"
	"Code_Review_N_1/internal/loader"
	"Code_Review_N_1/internal/repository"
	"Code_Review_N_1/internal/service"
	"github.com/gin-gonic/gin"
)

// ConfigDefaultInMemory is an struct that contains the configuration for the default application settings.
type ConfigDefaultInMemory struct {
	// FileLoader is the path to the file that contains the vehicles.
	FileLoader string
	// Addr is the address where the application will be listening.
	Addr string
}

// NewDefaultInMemory returns a new instance of a default application.
func NewDefaultInMemory(c *ConfigDefaultInMemory) *DefaultInMemory {
	// default config
	defaultCfg := &ConfigDefaultInMemory{
		FileLoader: "vehicles_100.json",
		Addr:       ":8080",
	}
	if c != nil {
		if c.FileLoader != "" {
			defaultCfg.FileLoader = c.FileLoader
		}
		if c.Addr != "" {
			defaultCfg.Addr = c.Addr
		}
	}

	return &DefaultInMemory{
		fileLoader: defaultCfg.FileLoader,
		addr:       defaultCfg.Addr,
	}
}

// DefaultInMemory is an struct that contains the default application settings.
type DefaultInMemory struct {
	// fileLoader is the path to the file that contains the vehicles.
	fileLoader string
	// addr is the address where the application will be listening.
	addr string
}

// Run starts the application.
func (d *DefaultInMemory) Run() (err error) {
	// dependencies initialization
	// loader
	ld := loader.NewVehicleJSON(d.fileLoader)
	data, err := ld.Load()
	if err != nil {
		return
	}

	// repository
	rp := repository.NewVehicleSlice(data.Data, data.LastId)

	// service
	sv := service.NewDefault(rp)

	// handler
	hd := handler.NewVehicleDefault(sv)

	// router
	rt := gin.New()
	// - middlewares
	rt.Use(gin.Logger())
	rt.Use(gin.Recovery())
	// - endpoints
	gr := rt.Group("/vehicles")
	{
		gr.GET("", hd.GetAll())
		gr.GET("/color/:color/year/:year", hd.FindByColorAndYear())
		gr.GET("/vehicles/brand/:brand/between/:start_year/:end_year", hd.FindByBrandAndYearRange())
		gr.GET("/average_speed/brand/:brand", hd.GetAverageSpeedByBrand())
		gr.GET("/fuel_type/:type", hd.GetByFuelType())
		gr.POST("", hd.AddVehicle())
		gr.POST("/batch", hd.AddMultipleVehicles())
		gr.PUT("/:id/update_speed", hd.UpdateMaxSpeed())
		gr.DELETE("/:id", hd.DeleteVehicle())

	}

	// run application
	err = rt.Run(d.addr)
	if err != nil {
		return
	}

	return
}
