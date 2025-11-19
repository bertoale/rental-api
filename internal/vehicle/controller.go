package vehicle

import (
	"go-rental/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller {
	return &Controller{
		service: s,
	}
}

func (ctrl *Controller) CreateVehicle(c *gin.Context) {
	var req VehicleRequest

	// Bind JSON
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}
	// Call service
	response, err := ctrl.service.CreateVehicle(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// Success
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "vehicle created successfully",
		"data":    response,
	})
}

func (ctrl *Controller) GetVehicles(c *gin.Context) {
	var filter VehicleFilter

	// Bind query parameters from URL
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid filter parameters: "+err.Error())
		return
	}

	vehicles, err := ctrl.service.GetAllVehicles(&filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "vehicles retrieved successfully", vehicles)
}

func (ctrl *Controller) GetVehicleByID(c *gin.Context) {
	id := c.Param("id")

	// Validate URL param
	vehicleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid vehicle ID format")
		return
	}

	// Service call
	vehicle, err := ctrl.service.GetVehicleByID(uint(vehicleID))
	if err != nil {
		if err.Error() == "vehicle not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Success response
	response.Success(c, http.StatusOK, "vehicle retrieved successfully", vehicle)
}

func (ctrl *Controller) UpdateVehicle(c *gin.Context) {
	id := c.Param("id")
	vehicleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid vehicle ID format")
		return
	}
	var req UpdateVehicleRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	updatedVehicle, err := ctrl.service.UpdateVehicle(uint(vehicleID), &req)
	if err != nil {
		if err.Error() == "vehicle not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "vehicle updated successfully", updatedVehicle)
}

func (ctrl *Controller) DeleteVehicle(c *gin.Context) {
	id := c.Param("id")
	vehicleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid vehicle ID format")
		return
	}
	err = ctrl.service.DeleteVehicle(uint(vehicleID))
	if err != nil {
		if err.Error() == "vehicle not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "vehicle deleted successfully", nil)
}
