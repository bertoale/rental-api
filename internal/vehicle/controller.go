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

// CreateVehicle godoc
// @Summary Create vehicle
// @Description Register a new vehicle
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body VehicleRequest true "Vehicle data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/vehicle/ [post]
func (ctrl *Controller) CreateVehicle(c *gin.Context) {
	var req VehicleRequest

	// Bind JSON
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	// Call service
	resp, err := ctrl.service.CreateVehicle(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// Success
	response.Success(c, http.StatusCreated, "vehicle created successfully", resp)
}

// GetVehicles godoc
// @Summary Get all vehicles
// @Description Retrieve all vehicles
// @Tags Vehicle
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/vehicle/ [get]
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

// GetVehicleByID godoc
// @Summary Get vehicle by ID
// @Description Retrieve a vehicle by its ID
// @Tags Vehicle
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/vehicle/{id} [get]
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

// UpdateVehicle godoc
// @Summary Update vehicle
// @Description Update vehicle data by ID
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vehicle ID"
// @Param data body UpdateVehicleRequest true "Vehicle update data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/vehicle/{id} [put]
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

// DeleteVehicle godoc
// @Summary Delete vehicle
// @Description Delete a vehicle by its ID
// @Tags Vehicle
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vehicle ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/vehicle/{id} [delete]
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
