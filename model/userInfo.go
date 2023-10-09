package model

type UserInfo struct {
	UserID       uint64 `gorm:"primary_key;NOT NULL" json:"userID"`
	Name         string `gorm:"NOT NULL" json:"name"`
	Sex          string `gorm:"NOT NULL" json:"sex"`
	Profession   string `gorm:"NOT NULL" json:"profession"`
	Position     string `gorm:"NOT NULL" json:"position"`
	Language     string `gorm:"NOT NULL" json:"language"`
	Domain       string `gorm:"NOT NULL" json:"domain"`
	Introduction string `gorm:"NOT NULL" json:"introduction"`
	Location     string `gorm:"NOT NULL" json:"location"`
	Email        string `gorm:"NOT NULL" json:"email"`
}
