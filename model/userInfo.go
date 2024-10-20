package model

type UserInfo struct {
	UserID       uint64 `gorm:"primary_key;NOT NULL"`
	Name         string
	Sex          string
	Profession   string
	Position     string
	Language     string
	Domain       string
	Introduction string
	Location     string
	Email        string
}
