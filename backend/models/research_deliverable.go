package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DeliverableProofCheckbox = "checkbox"
	DeliverableProofFile     = "file"

	DeliverableStatusPending = "pending"
	DeliverableStatusDone    = "done"
)

type ResearchDeliverable struct {
	InternalID            int64          `json:"internal_id" gorm:"primaryKey;column:deliverable_id"`
	PublicID              uuid.UUID      `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	ResearchInternalID    int64          `json:"-" gorm:"column:research_internal_id;not null;index"`
	Title                 string         `json:"title" gorm:"not null"`
	Description           string         `json:"description"`
	ProofType             string         `json:"proof_type" gorm:"column:proof_type;not null;default:'checkbox'"`
	Status                string         `json:"status" gorm:"not null;default:'pending'"`
	FilePath              *string        `json:"file_path,omitempty" gorm:"column:file_path"`
	CompletedByInternalID *int64         `json:"-" gorm:"column:completed_by_internal_id"`
	CompletedBy           *User          `json:"completed_by,omitempty" gorm:"foreignKey:CompletedByInternalID;references:InternalID"`
	CompletedAt           *time.Time     `json:"completed_at,omitempty" gorm:"column:completed_at"`
	DisplayOrder          int            `json:"display_order" gorm:"column:display_order;default:0"`
	CreatedByInternalID   int64          `json:"-" gorm:"column:created_by_internal_id;not null"`
	CreatedBy             *User          `json:"created_by,omitempty" gorm:"foreignKey:CreatedByInternalID;references:InternalID"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}
type DeliverableProgress struct {
	Total   int `json:"total"`
	Done    int `json:"done"`
	Pending int `json:"pending"`
	Percent int `json:"percent"`
}
