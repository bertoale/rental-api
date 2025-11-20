package customer

type Customer struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name"`
	Phone   string `json:"phone" gorm:"type:varchar(20);uniqueIndex"`
	Email   string `json:"email" gorm:"type:varchar(100);uniqueIndex"`
	Address string `json:"address"`
	IDCard  string `json:"id_card" gorm:"type:varchar(50);uniqueIndex"`
}

type CustomerRequest struct {
	Name    string `json:"name" form:"name" binding:"required"`
	Phone   string `json:"phone" form:"phone" binding:"required"`
	Email   string `json:"email" form:"email" binding:"required,email"`
	Address string `json:"address" form:"address" binding:"required"`
	IDCard  string `json:"id_card" form:"id_card" binding:"required"`
}

type CustomerResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
	IDCard  string `json:"id_card"`
}

type UpdateCustomerRequest struct {
	Name    *string `json:"name" form:"name" binding:"omitempty"`
	Phone   *string `json:"phone" form:"phone" binding:"omitempty"`
	Email   *string `json:"email" form:"email" binding:"omitempty,email"`
	Address *string `json:"address" form:"address" binding:"omitempty"`
	IDCard  *string `json:"id_card" form:"id_card" binding:"omitempty"`
}

type CustomerFilter struct {
	Name *string
}
