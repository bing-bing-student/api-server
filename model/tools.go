package model

type Tools struct {
	ToolID   uint64 `gorm:"primary_key;NOT NULL"`
	Describe string `gorm:"NOT NULL"`
	URL      string `gorm:"NOT NULL"`
}
