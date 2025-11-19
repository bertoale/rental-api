package vehicle

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(vehicle *Vehicle) error
	FindByID(id uint) (*Vehicle, error)
	FindAll(filter *VehicleFilter) ([]*Vehicle, error)
	Update(vehicle *Vehicle) error
	Delete(vehicle *Vehicle) error
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(vehicle *Vehicle) error {
	return r.db.Create(vehicle).Error
}

// Delete implements Repository.
func (r *repository) Delete(vehicle *Vehicle) error {
  return r.db.Delete(vehicle).Error
}


// FindAll implements Repository.
func (r *repository) FindAll(filter *VehicleFilter) ([]*Vehicle, error) {
    var vehicles []*Vehicle
    query := r.db.Model(&Vehicle{})
    // FILTER STATUS
    if filter.Status != nil {
        query = query.Where("status = ?", *filter.Status)
    }
    // FILTER BRAND
    if filter.Brand != nil {
        query = query.Where("brand LIKE ?", "%"+*filter.Brand+"%")
    }
    // FILTER MODEL
    if filter.Model != nil {
        query = query.Where("model LIKE ?", "%"+*filter.Model+"%")
    }
    // FILTER TYPE
    if filter.Type != nil {
        query = query.Where("type = ?", *filter.Type)
    }
    // FILTER YEAR RANGE
    if filter.MinYear != nil {
        query = query.Where("year >= ?", *filter.MinYear)
    }
    if filter.MaxYear != nil {
        query = query.Where("year <= ?", *filter.MaxYear)
    }

		if err := query.Find(&vehicles).Error; err != nil {
        return nil, err
    }

    return vehicles, nil
}


// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*Vehicle, error) {
	var v Vehicle
	if err := r.db.First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// Update implements Repository.
func (r *repository) Update(vehicle *Vehicle) error {
	return r.db.Save(vehicle).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}
