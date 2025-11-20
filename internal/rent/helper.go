package rent

func ToRentResponse(rent *Rent) *RentResponse {
	rentDate := ""
	if !rent.RentDate.IsZero() {
		rentDate = rent.RentDate.Format("2006-01-02 15:04:05")
	}

	returnDate := ""
	if rent.ReturnDate != nil && !rent.ReturnDate.IsZero() {
		returnDate = rent.ReturnDate.Format("2006-01-02 15:04:05")
	}

	return &RentResponse{
		ID:          rent.ID,
		Customer:    rent.Customer,
		Vehicle:     rent.Vehicle,
		RentDate:    rentDate,
		ReturnDate:  returnDate,
		TotalPrice:  rent.TotalPrice,
		Status:      rent.Status,
		Notes:       rent.Notes,
		CreatedBy:   rent.CreatedBy,
		UpdatedBy:   rent.UpdatedBy,
	}
}


// func calculateRentDays(start string, end string) (int, error) {
// 	layout := "2006-01-02" // format: YYYY-MM-DD
// 	t1, err := time.Parse(layout, start)
// 	if err != nil {
// 		return 0, errors.New("invalid rent_date format (use YYYY-MM-DD)")
// 	}
// 	t2, err := time.Parse(layout, end)
// 	if err != nil {
// 		return 0, errors.New("invalid return_date format (use YYYY-MM-DD)")
// 	}

// 	if t2.Before(t1) {
// 		return 0, errors.New("return date cannot be before rent date")
// 	}

// 	days := int(t2.Sub(t1).Hours()/24) + 1 // dihitung 1 hari penuh
// 	return days, nil
// }
