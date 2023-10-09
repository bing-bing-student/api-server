package model

type UserLogin struct {
	UserID   uint64 `gorm:"primary_key;NOT NULL" json:"userID"`
	Username string `gorm:"NOT NULL" json:"username"`
	Password string `gorm:"NOT NULL" json:"password"`
}
