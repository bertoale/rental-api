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

// Login godoc
// @Summary Login user
// @Description Login and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body LoginRequest true "Login data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/auth/login [post]
func (ctrl *Controller) Login (c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind((&req)); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
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
	response.Success(c, http.StatusOK, "login successful", gin.H{
		"user": user,
		"token": token,
	})
}

// CreateUser godoc
// @Summary Register user
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body RegisterRequest true "User data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/user/ [post]
func (ctrl *Controller) CreateUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := validator.FormatErrors(err)
		response.Error(c, http.StatusBadRequest, msg)
		return
	}

	result, err := ctrl.service.RegisterUser(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", result)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve all users
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/user/ [get]
func (ctrl *Controller) GetAllUsers(c *gin.Context) {
	u, err := ctrl.service.GetAllUsers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "user successfully retrieved", u)
}

// GetUsersByStatus godoc
// @Summary Get users by status
// @Description Retrieve users filtered by status
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param status query string true "User status"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/user/status [get]
func (ctrl *Controller) GetUsersByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		response.Error(c, http.StatusBadRequest, "status query is required")
		return
	}
	users, err := ctrl.service.GetUsersByStatus(status)
	if err != nil {
		// business error (misal: status tidak diizinkan)
		if err.Error() == "invalid status" {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// server/internal error
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "users successfully retrieved", users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/user/{id} [get]
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

// UpdateUser godoc
// @Summary Update user
// @Description Update user data by ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param data body UpdateRequest true "User update data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/user/{id} [put]
func (ctrl *Controller) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID")
		return
	}
	var req UpdateRequest
	if err := c.ShouldBind((&req)); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	updatedUser, err := ctrl.service.UpdateUser(uint(userID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "user successfully updated", updatedUser)
}