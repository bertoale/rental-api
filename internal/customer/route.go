package customer

import (
	"go-rental/pkg/config"
	"go-rental/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	customer := r.Group("/api/customer")
	{
		customer.POST("/", middlewares.Authenticate(cfg), ctrl.CreateCustomer)
		customer.GET("/", middlewares.Authenticate(cfg), ctrl.GetCustomers)
		customer.GET("/:id", middlewares.Authenticate(cfg), ctrl.GetCustomerByID)
		customer.PUT("/:id", middlewares.Authenticate(cfg), ctrl.UpdateCustomer)
	}
}