package user

import (
	"go-rental/pkg/config"
	"go-rental/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, ctrl *Controller, cfg *config.Config) {

	auth := r.Group("/api/auth")
{
		auth.POST("/login", ctrl.Login)
	}

	user := r.Group("/api/user")
	{
		user.POST("/",middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.CreateUser)
		user.GET("/status", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetUsersByStatus)
		user.GET("/", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetAllUsers)
		user.GET("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetUserByID)
		user.PUT("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.UpdateUser)
	}
}
