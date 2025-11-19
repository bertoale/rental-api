package vehicle

import (
	"go-rental/pkg/config"
	"go-rental/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupVehicleRoutes(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	vehicle := r.Group("/api/vehicle")
	{
		vehicle.POST("/", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.CreateVehicle)
		vehicle.GET("/", ctrl.GetVehicles)
		vehicle.GET("/:id", ctrl.GetVehicleByID)
		vehicle.PUT("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.UpdateVehicle)
		vehicle.DELETE("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.DeleteVehicle)
	}
}