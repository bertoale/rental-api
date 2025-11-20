package customer

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

func (ctrl *Controller) CreateCustomer(c *gin.Context) {
	var req CustomerRequest

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
	response, err := ctrl.service.CreateCustomer(&req)
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
		"message": "customer created successfully",
		"data":    response,
	})
}

func (ctrl *Controller) GetCustomers(c *gin.Context) {
	var filter CustomerFilter
	// Bind query parameters from URL
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid filter parameters: "+err.Error())
		return
	}

	customers, err := ctrl.service.GetAllCustomers(&filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "vehicles retrieved successfully", customers)

}

func (ctrl *Controller) GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	vehicleID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid customer ID: "+err.Error())
		return
	}

	customer, err := ctrl.service.GetCustomerByID(uint(vehicleID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "customer retrieved successfully", customer)

}

func (ctrl *Controller) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	vehicleID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid customer ID: "+err.Error())
		return
	}
	var req UpdateCustomerRequest
	// Bind JSON
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	customer, err := ctrl.service.UpdateCustomer(uint(vehicleID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "customer updated successfully", customer)
}
