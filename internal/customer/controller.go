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

// CreateCustomer godoc
// @Summary Create customer
// @Description Register a new customer
// @Tags Customer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body CustomerRequest true "Customer data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/customer/ [post]
func (ctrl *Controller) CreateCustomer(c *gin.Context) {
	var req CustomerRequest

	// Bind JSON
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	// Call service
	resp, err := ctrl.service.CreateCustomer(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// Success
	response.Success(c, http.StatusCreated, "customer created successfully", resp)
}

// GetCustomers godoc
// @Summary Get all customers
// @Description Retrieve all customers
// @Tags Customer
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/customer/ [get]
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

// GetCustomerByID godoc
// @Summary Get customer by ID
// @Description Retrieve a customer by their ID
// @Tags Customer
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/customer/{id} [get]
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

// UpdateCustomer godoc
// @Summary Update customer
// @Description Update customer data by ID
// @Tags Customer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param data body UpdateCustomerRequest true "Customer update data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/customer/{id} [put]
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
