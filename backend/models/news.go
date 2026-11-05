package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	NewsStatusDraft     = "draft"
	NewsStatusPublished = "published"
)

type News struct {
	InternalID         int64         `json:"-" gorm:"primaryKey;column:news_id"`
	PublicID           uuid.UUID     `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	Slug               string        `json:"slug" gorm:"column:slug;uniqueIndex;not null"`
	Title              string        `json:"title" gorm:"not null"`
	CategoryInternalID *int64        `json:"-" gorm:"column:category_internal_id;index"`
	Category           *NewsCategory `json:"category,omitempty" gorm:"foreignKey:CategoryInternalID;references:InternalID"`
	Content            string        `json:"content" gorm:"type:text;not null"`
	CoverImage         string        `json:"cover_image" gorm:"column:cover_image;not null"`
	Source             *string       `json:"source,omitempty" gorm:"column:source"`
	AuthorInternalID   int64         `json:"-" gorm:"column:author_internal_id;not null;index"`
	Author             *User         `json:"author,omitempty" gorm:"foreignKey:AuthorInternalID;references:InternalID"`
	Status             string        `json:"status" gorm:"not null;default:'draft'"`
	PublishedAt        *time.Time    `json:"published_at,omitempty" gorm:"column:published_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}