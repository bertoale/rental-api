package vehicle

import (
	"fmt"
	"go-rental/pkg/config"
)

type Service interface {
	CreateVehicle(req *VehicleRequest) (*VehicleResponse, error)
	GetVehicleByID(id uint) (*VehicleResponse, error)
	GetAllVehicles(filter *VehicleFilter) ([]*VehicleResponse, error)
	UpdateVehicle(id uint, req *UpdateVehicleRequest) (*VehicleResponse, error)
	DeleteVehicle(id uint) error
}

type service struct {
	repo Repository
}

// CreateVehicle implements Service.
func (s *service) CreateVehicle(req *VehicleRequest) (*VehicleResponse, error) {
	vehicle := &Vehicle{
		Type:        VehicleType(req.Type),
		PlateNumber: req.PlateNumber,
		Brand:       req.Brand,
		Model:       req.Model,
		Year:        req.Year,
		PricePerDay: req.PricePerDay,
		Status:      Avaibility(req.Status),
	}

	if err := s.repo.Create(vehicle); err != nil {
		return nil, err
	}

	return toVehicleResponse(vehicle), nil
}

// DeleteVehicle implements Service.
func (s *service) DeleteVehicle(id uint) error {
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}

	if err := s.repo.Delete(vehicle); err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	return nil
}

// GetAllVehicles implements Service.
func (s *service) GetAllVehicles(filter *VehicleFilter) ([]*VehicleResponse, error) {
	vehicles, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve vehicles: %w", err)
	}

	var responses []*VehicleResponse
	for _, v := range vehicles {
		responses = append(responses, toVehicleResponse(v))
	}

	return responses, nil
}

// GetVehicleByID implements Service.
func (s *service) GetVehicleByID(id uint) (*VehicleResponse, error) {
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("vehicle not found: %w", err)
	}
	return toVehicleResponse(vehicle), nil
}

// UpdateVehicle implements Service.
func (s *service) UpdateVehicle(id uint, req *UpdateVehicleRequest) (*VehicleResponse, error) {
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("vehicle not found: %w", err)
	}

	// Update only fields that are not nil
	if req.Type != nil {
		vehicle.Type = VehicleType(*req.Type)
	}
	if req.PlateNumber != nil {
		vehicle.PlateNumber = *req.PlateNumber
	}
	if req.Brand != nil {
		vehicle.Brand = *req.Brand
	}
	if req.Model != nil {
		vehicle.Model = *req.Model
	}
	if req.Year != nil {
		vehicle.Year = *req.Year
	}
	if req.PricePerDay != nil {
		vehicle.PricePerDay = *req.PricePerDay
	}
	if req.Status != nil {
		vehicle.Status = Avaibility(*req.Status)
	}

	// Save changes
	if err := s.repo.Update(vehicle); err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	return toVehicleResponse(vehicle), nil
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
	}
}
