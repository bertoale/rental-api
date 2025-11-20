package rent

import (
	"go-rental/internal/customer"
	"go-rental/internal/user"
	"go-rental/internal/vehicle"
	"time"
)

type RentStatus string

const (
	StatusOngoing   RentStatus = "ongoing"
	StatusCompleted RentStatus = "completed"
	StatusCancelled RentStatus = "cancelled"
)

type Rent struct {
	ID          uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID  uint        `json:"customer_id"`
	VehicleID   uint        `json:"vehicle_id"`
	RentDate    time.Time   `json:"rent_date"`
	ReturnDate  *time.Time  `json:"return_date" gorm:"default:null"`
	TotalPrice  float64     `json:"total_price"`
	Status      RentStatus  `json:"status"`
	Notes       string      `json:"notes"`

	CreatedByID uint `json:"created_by_id"`
	UpdatedByID uint `json:"updated_by_id"`

	// Relations
	CreatedBy user.User         `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedBy user.User         `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
	Customer  customer.Customer `json:"customer"   gorm:"foreignKey:CustomerID"`
	Vehicle   vehicle.Vehicle   `json:"vehicle"    gorm:"foreignKey:VehicleID"`
}

type RentRequest struct {
    CustomerID uint   `json:"customer_id" form:"customer_id" binding:"required"`
    VehicleID  uint   `json:"vehicle_id"  form:"vehicle_id"  binding:"required"`
    Notes      string `json:"notes"        form:"notes"`
}



type RentResponse struct {
	ID          uint        			`json:"id"`
	Customer    customer.Customer `json:"customer"`
	Vehicle     vehicle.Vehicle   `json:"vehicle"`
	RentDate    string      			`json:"rent_date"`
	ReturnDate  string      			`json:"return_date"`
	TotalPrice  float64    				`json:"total_price"`
	Status      RentStatus 				`json:"status"`
	Notes       string     				`json:"notes"`
	CreatedBy   user.User         `json:"created_by"`
	UpdatedBy   user.User         `json:"updated_by"`
}

type UpdateRentRequest struct {
    Status *string `json:"status" form:"status" binding:"omitempty"`
    Notes  *string `json:"notes"  form:"notes"  binding:"omitempty"`
}
