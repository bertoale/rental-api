package user

type RoleType string
type StatusType string

const (
	RoleAdmin  RoleType = "admin"
	RoleMember RoleType = "staff"
)

const (
	StatusActive   StatusType = "active"
	StatusInactive StatusType = "inactive"
)

type User struct {
	ID       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string     `json:"name"`
	Phone    string     `json:"phone" gorm:"uniqueIndex"`
	Username string     `json:"username" gorm:"uniqueIndex"`
	Password string     `json:"-"`
	Role     RoleType   `json:"role" gorm:"type:enum('admin', 'staff');default:'staff'"`
	Status   StatusType `json:"status" gorm:"type:enum('active', 'inactive');default:'active'"`
}

type RegisterRequest struct {
	Name     string   `json:"name" binding:"required"`
	Phone    string   `json:"phone" binding:"required,e164"`
	Username string   `json:"username" binding:"required,alphanum"`
	Password string   `json:"password" binding:"required,min=8"`
	Role     RoleType `json:"role" binding:"required,oneof=admin staff"`
}

type RegisterResponse struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}
