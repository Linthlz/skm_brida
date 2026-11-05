package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewsCategory struct {
	InternalID   int64          `json:"-" gorm:"primaryKey;column:category_id"`
	PublicID     uuid.UUID      `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	Name         string         `json:"name" gorm:"not null;uniqueIndex"`
	Slug         string         `json:"slug" gorm:"not null;uniqueIndex"`
	DisplayOrder int            `json:"-" gorm:"column:display_order;default:0"`
	IsActive     bool           `json:"-" gorm:"column:is_active;default:true"`
	NewsList     []News         `json:"-" gorm:"foreignKey:CategoryInternalID;references:InternalID"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
