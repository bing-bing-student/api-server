package model

type Label struct {
	LabelID   uint64 `gorm:"primary_key;NOT NULL" json:"labelID"`
	LabelName string `gorm:"NOT NULL" json:"labelName"`
}
