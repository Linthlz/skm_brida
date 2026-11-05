package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupDocumentation struct {
	InternalID       int64                `json:"internal_id" gorm:"primaryKey;column:documentation_id"`
	PublicID         uuid.UUID            `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	GroupInternalID  int64                `json:"-" gorm:"column:group_internal_id;not null;index"`
	AuthorInternalID int64                `json:"-" gorm:"column:author_internal_id;not null;index"`
	Author           *User                `json:"author,omitempty" gorm:"foreignKey:AuthorInternalID;references:InternalID"`
	Title            string               `json:"title" gorm:"not null"`
	Narrative        string               `json:"narrative"`
	Location         string               `json:"location"`
	VideoURL         *string              `json:"video_url,omitempty" gorm:"column:video_url"`
	EventDate        *time.Time           `json:"event_date,omitempty" gorm:"column:event_date;type:date"`
	Images           []DocumentationImage `json:"images,omitempty" gorm:"foreignKey:DocumentationInternalID;references:InternalID"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	DeletedAt        gorm.DeletedAt       `json:"-" gorm:"index"`
}

type DocumentationImage struct {
	InternalID              int64          `json:"internal_id" gorm:"primaryKey;column:image_id"`
	PublicID                uuid.UUID      `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	DocumentationInternalID int64          `json:"-" gorm:"column:documentation_internal_id;not null;index"`
	FilePath                string         `json:"file_path" gorm:"column:file_path;not null"`
	Caption                 string         `json:"caption"`
	DisplayOrder            int            `json:"display_order" gorm:"column:display_order;default:0"`
	CreatedAt               time.Time      `json:"created_at"`
	DeletedAt               gorm.DeletedAt `json:"-" gorm:"index"`
}
