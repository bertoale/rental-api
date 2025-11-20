package rent

import (
	"go-rental/internal/customer"
	"go-rental/internal/user"
	"go-rental/internal/vehicle"

	"gorm.io/gorm"
)

type Repository interface {
	Create(rent *Rent) error
	FindByID(id uint) (*Rent, error)
	FindAll() ([]*Rent, error)
	Update(rent *Rent) error
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(rent *Rent) error {
	return r.db.Create(rent).Error
}

// Delete implements Repository.


// FindAll implements Repository.
func (r *repository) FindAll() ([]*Rent, error) {
	var rents []*Rent
	if err := r.db.Preload("CreatedBy").Preload("UpdatedBy").Preload("Vehicle").Preload("Customer").Find(&rents).Error; err != nil {
		return nil, err
	}
	return rents, nil
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*Rent, error) {
	var rent Rent
	if err := r.db.Preload("CreatedBy").Preload("UpdatedBy").Preload("Vehicle").Preload("Customer").First(&rent, id).Error; err != nil {
		return nil, err
	}
	return &rent, nil
}

// Update implements Repository.
func (r *repository) Update(rent *Rent) error {
	return r.db.Save(rent).Error
}

func NewRepository(db *gorm.DB, userRepo user.Repository, vehicleRepo vehicle.Repository, customerRepo customer.Repository) Repository {
	return &repository{
		db: db,
	}
}
