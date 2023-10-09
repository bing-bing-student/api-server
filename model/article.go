package model

import (
	"time"
)

type Article struct {
	ArticleID   uint64    `gorm:"primary_key;NOT NULL" json:"articleID"`
	ViewNum     uint64    `gorm:"NOT NULL" json:"viewNum"`
	BlogTitle   string    `gorm:"NOT NULL" json:"blogTitle"`
	ShowContext string    `gorm:"NOT NULL" json:"showContext"`
	Context     string    `gorm:"type:text;NOT NULL" json:"context"`
	CreatedTime time.Time `gorm:"NOT NULL" json:"createdAt"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updatedAt"`
	DeletedTime time.Time `gorm:"column:deleted_time" json:"deletedAt"`
	UserID      uint64    `gorm:"NOT NULL" json:"userID"`
	LabelID     uint64    `gorm:"NOT NULL" json:"labelID"`
}
