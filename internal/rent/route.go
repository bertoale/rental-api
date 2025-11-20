package rent

import (
	"go-rental/pkg/config"
	"go-rental/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func RentSetupRoutes(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	rent := r.Group("/api/rent")
	{
		rent.POST("/", middlewares.Authenticate(cfg), ctrl.CreateRent)
		rent.GET("/", middlewares.Authenticate(cfg), ctrl.GetRents)
		rent.GET("/:id", middlewares.Authenticate(cfg), ctrl.GetRentByID)
		rent.PUT("/:id/", middlewares.Authenticate(cfg), ctrl.UpdateRent)
	}
}