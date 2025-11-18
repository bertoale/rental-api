package user

import "gorm.io/gorm"

type Repository interface {
	//auth
	FindByUsername(username string) (*User, error)
	//user
	CreateUser(user *User) error
	FindByID(id uint) (*User, error)
	FindByStatus(status StatusType) ([]*User, error)
	FindAll() ([]*User, error)
	Update(user *User) error
}

type repository struct {
	db *gorm.DB
}

// FindByStatus implements Repository.
func (r *repository) FindByStatus(status StatusType) ([]*User, error) {
	var u []*User
	if err := r.db.Where("status = ?", status).Order("name desc").Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// CreateUser implements Repository.
func (r *repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}



// FindAll implements Repository.
func (r *repository) FindAll() ([]*User, error) {
	var u []*User
	if err := r.db.Order("name desc").Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByUsername implements Repository.
func (r *repository) FindByUsername(username string) (*User, error) {
	var u User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Update implements Repository.
func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}
