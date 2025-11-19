package vehicle

func toVehicleResponse(v *Vehicle) *VehicleResponse {
	return &VehicleResponse{
		ID:          v.ID,
		Type:        v.Type,
		PlateNumber: v.PlateNumber,
		Brand:       v.Brand,
		Model:       v.Model,
		Year:        v.Year,
		PricePerDay: v.PricePerDay,
		Status:      v.Status,
	}
}
