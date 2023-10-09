package model

type Tools struct {
	ToolID   uint64 `gorm:"primary_key;NOT NULL" json:"toolID"`
	Describe string `gorm:"NOT NULL" json:"describe"`
	URL      string `gorm:"NOT NULL" json:"url"`
}
