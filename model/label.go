package model

type Label struct {
	LabelID   uint64 `gorm:"primary_key;NOT NULL"`
	LabelName string `gorm:"NOT NULL"`
}
