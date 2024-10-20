package model

type UserLogin struct {
	UserID   uint64 `gorm:"primary_key;NOT NULL"`
	Username string `gorm:"NOT NULL"`
	Password string `gorm:"NOT NULL"`
	Phone    string `gorm:"NOT NULL"`
}
