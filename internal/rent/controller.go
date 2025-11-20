package rent

import (
	"go-rental/internal/customer"
	"go-rental/internal/vehicle"
	"go-rental/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	rentService     Service
	vehicleService  vehicle.Service
	customerService customer.Service
}

func NewController(rentService Service, vehicleService vehicle.Service, customerService customer.Service) *Controller {
	return &Controller{
		rentService:     rentService,
		vehicleService:  vehicleService,
		customerService: customerService,
	}
}

func (ctrl *Controller) CreateRent(c *gin.Context) {
	var req RentRequest

	// Bind JSON request
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// Validate customer exists
	_, err := ctrl.customerService.GetCustomerByID(req.CustomerID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "customer not found")
		return
	}

	// Validate vehicle exists
	_, err = ctrl.vehicleService.GetVehicleByID(req.VehicleID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "vehicle not found")
		return
	}

	// Call service to create rent
	rent, err := ctrl.rentService.CreateRent(&req, userID.(uint))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Success response
	response.Success(c, http.StatusCreated, "rent created successfully", rent)
}

func (ctrl *Controller) GetRents(c *gin.Context) {
	rents, err := ctrl.rentService.GetAllRents()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "rents retrieved successfully", rents)
}

func (ctrl *Controller) GetRentByID(c *gin.Context) {
	id := c.Param("id")
	rentID,err := strconv.ParseUint(id,10,32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid rent ID")
		return
	}
	rent, err := ctrl.rentService.GetRentByID(uint(rentID))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "rent retrieved successfully", rent)
}

func (ctrl *Controller) UpdateRent(c *gin.Context) {
    // Parse rent ID
    idParam := c.Param("id")
    rentID, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "invalid rent ID")
        return
    }

    // Bind request
    var req UpdateRentRequest
    if err := c.ShouldBind(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "invalid request body: "+err.Error())
        return
    }

    // Get user ID from context
    userID, exists := c.Get("userID")
    if !exists {
        response.Error(c, http.StatusUnauthorized, "user not authenticated")
        return
    }

    // Call service
    updatedRent, err := ctrl.rentService.UpdateRent(uint(rentID), &req, userID.(uint))
    if err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    // Success
    response.Success(c, http.StatusOK, "rent updated successfully", updatedRent)
}
