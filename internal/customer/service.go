package customer

import (
	"fmt"
	"go-rental/pkg/config"
)

type Service interface {
	CreateCustomer(req *CustomerRequest) (*CustomerResponse, error)
	GetAllCustomers(filter *CustomerFilter) ([]*CustomerResponse, error)
	UpdateCustomer(id uint, req *UpdateCustomerRequest) (*CustomerResponse, error)
	GetCustomerByID(id uint) (*CustomerResponse, error)
}

type service struct {
	repo Repository
}

// CreateCustomer implements Service.
func (s *service) CreateCustomer(req *CustomerRequest) (*CustomerResponse, error) {
	customer := &Customer{
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
		IDCard:  req.IDCard,
	}
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}
	return ToCustomerResponse(customer), nil

}

// GetAllCustomers implements Service.
func (s *service) GetAllCustomers(filter *CustomerFilter) ([]*CustomerResponse, error) {
	customers, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customer: %w", err)
	}
	var responses []*CustomerResponse
	for _, customer := range customers {
		responses = append(responses, ToCustomerResponse(customer))
	}
	return responses, nil
}

// GetCustomerByID implements Service.
func (s *service) GetCustomerByID(id uint) (*CustomerResponse, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	return ToCustomerResponse(customer), nil

}

// UpdateCustomer implements Service.
func (s *service) UpdateCustomer(id uint, req *UpdateCustomerRequest) (*CustomerResponse, error) {
	
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.Email != nil {
		customer.Email = *req.Email
	}
	if req.Address != nil {
		customer.Address = *req.Address
	}
	if req.IDCard != nil {
		customer.IDCard = *req.IDCard
	}
	if err := s.repo.Update(customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}
	
	return ToCustomerResponse(customer), nil
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{repo: repo}
}
