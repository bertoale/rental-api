package vehicle

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
	ID          uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Type        VehicleType `json:"type" gorm:"type:enum('car', 'bike')"`
	PlateNumber string      `json:"plate_number" gorm:"uniqueIndex"`
	Brand       string      `json:"brand"`
	Model       string      `json:"model"`
	Year        int         `json:"year"`
	PricePerDay float64     `json:"price_per_day"`
	Status      Avaibility  `json:"status" gorm:"type:enum('available', 'rented', 'maintenance');default:'available'"`
}

type VehicleRequest struct {
	Type        string `json:"type" form:"type" binding:"required"`
	PlateNumber string `json:"plate_number" form:"plate_number" binding:"required"`
	Brand       string `json:"brand" form:"brand" binding:"required"`
	Model       string
}