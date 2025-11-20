package customer

func ToCustomerResponse(customer *Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:      customer.ID,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Email:   customer.Email,
		Address: customer.Address,
		IDCard:  customer.IDCard,
	}
}