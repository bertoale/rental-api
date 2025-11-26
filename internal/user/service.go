package user

import (
	"errors"
	"go-rental/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Service interface {
	// Auth
	Login(req LoginRequest) (string, *UserResponse, error)
	GenerateToken(user *User) (string, error)

	// User
	RegisterUser(req RegisterRequest) (*UserResponse, error)
	GetUserByID(id uint) (*UserResponse, error)
	GetUsersByStatus(status string) ([]UserResponse, error)
	GetAllUsers() ([]UserResponse, error)
	UpdateUser(userID uint, req *UpdateRequest) (*UserResponse, error)
}

type service struct {
	repo Repository
	cfg  *config.Config
}



func (s *service) GenerateToken(user *User) (string, error) {
	duration, err := time.ParseDuration(s.cfg.JWTExpires)
	if err != nil {
		duration = 168 * time.Hour // default 7 days
	}

	claims := Claims{
		ID:   user.ID,
		Role: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.JWTSecret))
}



func (s *service) GetAllUsers() ([]UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, errors.New("failed to retrieve users")
	}

	var responses []UserResponse
	for _, user := range users {
		responses = append(responses, UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Phone:    user.Phone,
			Role:     string(user.Role),
			Status:   string(user.Status),
		})
	}

	return responses, nil
}



func (s *service) GetUserByID(id uint) (*UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return ToUserResponse(u), nil
}


func (s *service) GetUsersByStatus(status string) ([]UserResponse, error) {
	users, err := s.repo.FindByStatus(StatusType(status))
	if err != nil {
		return nil, errors.New("failed to retrieve users by status")
	}

	var responses []UserResponse
	for _, user := range users {
		responses = append(responses, *ToUserResponse(user))
	}

	return responses, nil
}

func (s *service) Login(req LoginRequest) (string, *UserResponse, error) {
	// required validation handled by Gin

	u, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	token, err := s.GenerateToken(u)
	if err != nil {
		return "", nil, err
	}

	return token, ToUserResponse(u), nil
}


func (s *service) RegisterUser(req RegisterRequest) (*UserResponse, error) {
	// validation required handled by Gin

	existingUser, err := s.repo.FindByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check username")
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Name:     req.Name,
		Phone:    req.Phone,
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     req.Role,
		Status:   StatusActive,
	}

	if err := s.repo.CreateUser(u); err != nil {
		return nil, err
	}

	return ToUserResponse(u), nil
}



func (s *service) UpdateUser(userID uint, req *UpdateRequest) (*UserResponse, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Name
	if req.Name != nil && *req.Name != "" {
		u.Name = *req.Name
	}

	// Phone
	if req.Phone != nil && *req.Phone != "" {
		u.Phone = *req.Phone
	}

	// Username
	if req.Username != nil && *req.Username != "" {
		u.Username = *req.Username
	}

	// Role
	if req.Role != nil && *req.Role != "" {
		u.Role = RoleType(*req.Role)
	}

	// Status
	if req.Status != nil && *req.Status != "" {
		u.Status = StatusType(*req.Status)
	}

	// Password (hash only if not empty)
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.Password = string(hashedPassword)
	}

	// Save changes
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}

	return ToUserResponse(u), nil
}



func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}
