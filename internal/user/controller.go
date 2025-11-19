package user

import (
	"go-rental/pkg/config"
	"go-rental/pkg/response"
	"go-rental/pkg/validator"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
	cfg		*config.Config
}

func NewController(s Service, cfg *config.Config) *Controller {
	return &Controller{
		service: s,
		cfg:     cfg,
	}
}

func (ctrl *Controller) Login (c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind((&req)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, user, err := ctrl.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie(
		"token",
		token,
		int((7*24*time.Hour).Seconds()),
		"/",
		"",
		ctrl.cfg.NodeEnv == "production",
		true,
	)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"token": token})
}

func (ctrl *Controller) CreateUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := validator.FormatErrors(err)
		response.Error(c, 400, msg)
		return
	}

	result, err := ctrl.service.RegisterUser(req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 201, "User created successfully", result)
}



func (ctrl *Controller) GetAllUsers(c *gin.Context) {
	u, err := ctrl.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user successfully retrieved",
		"users": u})
}

func (ctrl *Controller) GetUsersByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "status query is required",
		})
		return
	}
	users, err := ctrl.service.GetUsersByStatus(status)
	if err != nil {
		// business error (misal: status tidak diizinkan)
		if err.Error() == "invalid status" {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		// server/internal error
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "users successfully retrieved",
		"users":   users,
	})
}

func (ctrl *Controller) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user ID"})
		return
	}
	user, err := ctrl.service.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user successfully retrieved",
		"user":    user,
	})
}

func (ctrl *Controller) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user ID"})
		return
	}
	var req UpdateRequest
	if err := c.ShouldBind((&req)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	updatedUser, err := ctrl.service.UpdateUser(uint(userID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    updatedUser,
	})
}