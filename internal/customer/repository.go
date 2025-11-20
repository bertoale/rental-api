package customer

import "gorm.io/gorm"

type Repository interface {
	Create(customer *Customer) error
	FindByID(id uint) (*Customer, error)
	FindAll(filter *CustomerFilter) ([]*Customer, error)
	Update(customer *Customer) error
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(customer *Customer) error {
	return r.db.Create(customer).Error
}

// FindAll implements Repository.
func (r *repository) FindAll(filter *CustomerFilter) ([]*Customer, error) {
	var customers []*Customer
	query := r.db.Model(&Customer{})
	// FILTER NAME
	if filter.Name != nil {
		query = query.Where("name LIKE ?", "%"+*filter.Name+"%")
	}
	if err := query.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*Customer, error) {
	var customer Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// Update implements Repository.
func (r *repository) Update(customer *Customer) error {
	return r.db.Save(customer).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
