package model

import (
	"time"
)

type Article struct {
	ArticleID   uint64    `gorm:"primary_key;NOT NULL"`
	ViewNum     uint64    `gorm:"NOT NULL"`
	BlogTitle   string    `gorm:"NOT NULL"`
	ShowContext string    `gorm:"NOT NULL"`
	Context     string    `gorm:"type:text;NOT NULL"`
	CreatedTime time.Time `gorm:"NOT NULL"`
	UpdatedTime time.Time `gorm:"NOT NULL"`
	LabelID     uint64    `gorm:"NOT NULL"`
}
