package rent

import (
	"errors"
	"go-rental/internal/vehicle"
	"go-rental/pkg/config"
	"time"
)

type Service interface {
	CreateRent(req *RentRequest, createdBy uint) (*RentResponse, error)
	GetRentByID(id uint) (*RentResponse, error)
	GetAllRents() ([]*RentResponse, error)
	UpdateRent(id uint, req *UpdateRentRequest, updatedBy uint) (*RentResponse, error)
}

type service struct {
	vehicleRepo vehicle.Repository
	repo        Repository
    cfg         config.Config
}

// CreateRent implements Service.
func (s *service) CreateRent(req *RentRequest, createdBy uint) (*RentResponse, error) {
    // 1. Cek vehicle
    vh, err := s.vehicleRepo.FindByID(req.VehicleID)
    if err != nil {
        return nil, errors.New("vehicle not found")
    }
    if vh.Status != vehicle.StatusAvailable {
        return nil, errors.New("vehicle is not available")
    }    // 2. Buat rent dengan RentDate otomatis (sekarang)
    rent := &Rent{
        CustomerID:  req.CustomerID,
        VehicleID:   req.VehicleID,
        RentDate:    time.Now(), // Set otomatis saat dibuat
        Status:      StatusOngoing,
        Notes:       req.Notes,
        TotalPrice:  0, // Akan dihitung saat completed
        CreatedByID: createdBy,
        UpdatedByID: createdBy,
    }

    if err := s.repo.Create(rent); err != nil {
        return nil, err
    }

    // 3. Update status kendaraan
    vh.Status = vehicle.StatusRented
    if err := s.vehicleRepo.Update(vh); err != nil {
        return nil, errors.New("failed to update vehicle status")
    }

    // 4. Load relasi (customer, vehicle, created_by, updated_by)
    createdRent, err := s.repo.FindByID(rent.ID)
    if err != nil {
        return nil, err
    }

    return ToRentResponse(createdRent), nil
}



// GetAllRents implements Service.
func (s *service) GetAllRents() ([]*RentResponse, error) {
	rents, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var responses []*RentResponse
	for _, rent := range rents {
		responses = append(responses, ToRentResponse(rent))
	}
	return responses, nil
}

// GetRentByID implements Service.
func (s *service) GetRentByID(id uint) (*RentResponse, error) {
	rent, err := s.repo.FindByID(id)	
	if err != nil {
		return nil, err
	}	
	return ToRentResponse(rent), nil
}

// UpdateRent implements Service.
func (s *service) UpdateRent(id uint, req *UpdateRentRequest, updatedBy uint) (*RentResponse, error) {
    rent, err := s.repo.FindByID(id)
    if err != nil {
        return nil, errors.New("rent not found")
    }

    // Update UpdatedByID
    rent.UpdatedByID = updatedBy    // Update Notes
    if req.Notes != nil {
        rent.Notes = *req.Notes
    }

    // Update Status
    if req.Status != nil {
        oldStatus := rent.Status
        newStatus := RentStatus(*req.Status)

        // Validasi status transition
        if newStatus != StatusOngoing && newStatus != StatusCompleted && newStatus != StatusCancelled {
            return nil, errors.New("invalid status value")
        }

        // Validasi: tidak bisa complete/cancel jika sudah complete
        if oldStatus == StatusCompleted && newStatus != StatusCompleted {
            return nil, errors.New("cannot change status from completed")
        }

        // Validasi: tidak bisa complete/cancel jika sudah cancelled
        if oldStatus == StatusCancelled && newStatus != StatusCancelled {
            return nil, errors.New("cannot change status from cancelled")
        }

        rent.Status = newStatus

        // Jika status berubah menjadi completed
        if newStatus == StatusCompleted && oldStatus != StatusCompleted {
            // Set return date ke sekarang
            now := time.Now()
            rent.ReturnDate = &now

            // Hitung total price
            vh, err := s.vehicleRepo.FindByID(rent.VehicleID)
            if err != nil {
                return nil, errors.New("vehicle not found")
            }

            // Hitung jumlah hari
            days := int(rent.ReturnDate.Sub(rent.RentDate).Hours()/24) + 1
            if days < 1 {
                days = 1
            }
            rent.TotalPrice = float64(days) * vh.PricePerDay

            // Update status kendaraan menjadi available
            vh.Status = vehicle.StatusAvailable
            if err := s.vehicleRepo.Update(vh); err != nil {
                return nil, errors.New("failed to update vehicle status")
            }
        }

        // Jika status berubah menjadi cancelled
        if newStatus == StatusCancelled && oldStatus != StatusCancelled {
            // Update status kendaraan menjadi available
            vh, err := s.vehicleRepo.FindByID(rent.VehicleID)
            if err != nil {
                return nil, errors.New("vehicle not found")
            }
            vh.Status = vehicle.StatusAvailable
            if err := s.vehicleRepo.Update(vh); err != nil {
                return nil, errors.New("failed to update vehicle status")
            }
        }
    }

    // Save to DB
    if err := s.repo.Update(rent); err != nil {
        return nil, errors.New("failed to update rent")
    }

    // Reload relasi agar response lengkap
    updatedRent, err := s.repo.FindByID(rent.ID)
    if err != nil {
        return nil, err
    }
    return ToRentResponse(updatedRent), nil
}

func NewService(repo Repository, vehicleRepo vehicle.Repository, cfg config.Config) Service {
    return &service{
        repo:        repo,
        vehicleRepo: vehicleRepo,
        cfg:         cfg,
    }
}