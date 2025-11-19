package vehicle

import "gorm.io/gorm"

type VehicleType string
type Avaibility string

const (
	VehicleCar  VehicleType = "car"
	VehicleBike VehicleType = "bike"
)
const (
	StatusAvailable   Avaibility = "available"
	StatusRented      Avaibility = "rented"
	StatusMaintenance Avaibility = "maintenance"
)

type Vehicle struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Type        VehicleType    `json:"type" gorm:"type:enum('car', 'bike')"`
	PlateNumber string         `json:"plate_number" gorm:"type:varchar(20);uniqueIndex"`
	Brand       string         `json:"brand"`
	Model       string         `json:"model"`
	Year        int            `json:"year"`
	PricePerDay float64        `json:"price_per_day"`
	Status      Avaibility     `json:"status" gorm:"type:enum('available', 'rented', 'maintenance');default:'available'"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type VehicleRequest struct {
	Type        string  `json:"type" form:"type" binding:"required"`
	PlateNumber string  `json:"plate_number" form:"plate_number" binding:"required"`
	Brand       string  `json:"brand" form:"brand" binding:"required"`
	Model       string  `json:"model" form:"model" binding:"required"`
	Year        int     `json:"year" form:"year" binding:"required"`
	PricePerDay float64 `json:"price_per_day" form:"price_per_day" binding:"required"`
	Status      string  `json:"status" form:"status" binding:"required"`
}

type VehicleResponse struct {
	ID          uint        `json:"id"`
	Type        VehicleType `json:"type"`
	PlateNumber string      `json:"plate_number"`
	Brand       string      `json:"brand"`
	Model       string      `json:"model"`
	Year        int         `json:"year"`
	PricePerDay float64     `json:"price_per_day"`
	Status      Avaibility  `json:"status"`
}

type UpdateVehicleRequest struct {
	Type        *string  `json:"type" form:"type" binding:"omitempty"`
	PlateNumber *string  `json:"plate_number" form:"plate_number" binding:"omitempty"`
	Brand       *string  `json:"brand" form:"brand" binding:"omitempty"`
	Model       *string  `json:"model" form:"model" binding:"omitempty"`
	Year        *int     `json:"year" form:"year" binding:"omitempty"`
	PricePerDay *float64 `json:"price_per_day" form:"price_per_day" binding:"omitempty"`
	Status      *string  `json:"status" form:"status" binding:"omitempty"`
}

type VehicleFilter struct {
    Status *string
    Brand  *string
    Model  *string
    Type   *string
    MinYear *int
    MaxYear *int
}
