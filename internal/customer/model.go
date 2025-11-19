package customer

type Customer struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name"`
	Phone   string `json:"phone" gorm:"type:varchar(20);uniqueIndex"`
	Email   string `json:"email" gorm:"uniqueIndex"`
	Address string `json:"address"`
	IDCard  string `json:"id_card" gorm:"uniqueIndex"`
}