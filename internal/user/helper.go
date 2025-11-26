package user

func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Phone:    user.Phone,
		Role:     string(user.Role),
		Status:   string(user.Status),
	}
}