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
	//auth
	Login(req LoginRequest) (string, *UserResponse, error)
	GenerateToken(user *User) (string, error)
	//user
	RegisterUser(req RegisterRequest) (*UserResponse,error)
	GetUserByID(id uint) (*UserResponse, error)
	GetUsersByStatus(status string) ([]UserResponse, error)
	GetAllUsers() ([]UserResponse, error)
	UpdateUser(userID uint, req *UpdateRequest) (*UserResponse, error)
}

type service struct {
	repo Repository
	cfg  *config.Config
}

// GenerateToken implements Service.
func (s *service) GenerateToken(user *User) (string, error) {
	// Parse duration from config
	duration, err := time.ParseDuration(s.cfg.JWTExpires)
	if (err != nil) {
		duration = 168 * time.Hour // Default 7 days
	}
	// Create claims with user ID and standard claims
	claims := Claims{
		ID: user.ID,
		Role: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// Create token with signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token with secret key and return token string
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// DeleteUser implements Service.


// GetAllUsers implements Service.
func (s *service) GetAllUsers() ([]UserResponse, error) {
	u , err := s.repo.FindAll()
	if err != nil {
		return nil, errors.New("failed to retrieve users")
	}
	
	var responses []UserResponse
	for _, user := range u {
	response := UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Phone:    user.Phone,
		Role:     string(user.Role),
		Status:   string(user.Status),
	}
	responses = append(responses, response)
	}
	return responses, nil
}
// GetUserByID implements Service.
func (s *service) GetUserByID(id uint) (*UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := &UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Phone:    u.Phone,
		Role:     string(u.Role),
		Status:   string(u.Status),
	}
	return response, nil
}

// GetUsersByStatus implements Service.
func (s *service) GetUsersByStatus(status string) ([]UserResponse, error) {
	u, err := s.repo.FindByStatus(StatusType(status))
	if err != nil {
		return nil, errors.New("failed to retrieve users by status")
	}
	var responses []UserResponse
	for _, user := range u {
		response := UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Phone:    user.Phone,
			Role:     string(user.Role),
			Status:   string(user.Status),
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// Login implements Service.
func (s *service) Login(req LoginRequest) (string, *UserResponse, error) {
	if req.Username == "" || req.Password == "" {
		return "", nil, errors.New("username and password are required")
	}

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
	userResp := &UserResponse{
		ID:   u.ID,
		Name: u.Name,
		Username: u.Username,
		Phone: u.Phone,
		Role: string(u.Role),
		Status: string(u.Status),
	}
	return token, userResp, nil
	
}

// RegisterUser implements Service.
func (s *service) RegisterUser(req RegisterRequest) (*UserResponse,error) {
	if req.Name == "" || req.Phone == "" || req.Username == "" || req.Password == "" {
			return nil, errors.New("all fields are required")
	}
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
	if err := s.repo.CreateUser(u);  err != nil {
		return nil, err
	}
	
	return &UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Phone:    u.Phone,
		Role:     string(u.Role),
		Status:   string(u.Status),
	}, nil
}

// UpdateUser implements Service.
func (s *service)	UpdateUser(userID uint, req *UpdateRequest) (*UserResponse, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("name cannot be empty")
		}
		u.Name = *req.Name
	}

	if req.Phone != nil {
		if *req.Phone == "" {
			return nil, errors.New("phone cannot be empty")
		}
		u.Phone = *req.Phone
	}

	if req.Username != nil {
		if *req.Username == "" {
			return nil, errors.New("username cannot be empty")
		}
		u.Username = *req.Username
	}

	if req.Role != nil {
		if *req.Role == "" {
			return nil, errors.New("role cannot be empty")
		}
		u.Role = RoleType(*req.Role)
	}

	if req.Status != nil {
		if *req.Status == "" {
			return nil, errors.New("status cannot be empty")
		}
		u.Status = StatusType(*req.Status)
	}

	if req.Password != nil {
		if *req.Password == "" {
			return nil, errors.New("password cannot be empty")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.Password = string(hashedPassword)
	}

	if err := s.repo.Update(u); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Phone:    u.Phone,
		Role:     string(u.Role),
		Status:   string(u.Status),
	}, nil
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}
