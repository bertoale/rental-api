package rent

import (
	"go-rental/internal/customer"
	"go-rental/internal/user"
	"go-rental/internal/vehicle"
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
	RentDate    string      `json:"rent_date"`
	ReturnDate  string      `json:"return_date"`
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
type CreateRentRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	VehicleID  uint   `json:"vehicle_id" binding:"required"`
	RentDate   string `json:"rent_date" binding:"required,datetime=2006-01-02"`
	ReturnDate string `json:"return_date" binding:"required,datetime=2006-01-02,gtfield=RentDate"`
	Notes      string `json:"notes"`
}