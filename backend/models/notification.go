package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	InternalID          int64      `json:"internal_id" gorm:"primaryKey;column:notification_id"`
	PublicID            uuid.UUID  `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	RecipientInternalID int64      `json:"-" gorm:"column:recipient_internal_id;not null;index"`
	Subject             string     `json:"subject"`
	Message             string     `json:"message" gorm:"type:text;not null"`
	Type                string     `json:"type" gorm:"column:type;default:'info'"`
	ReferenceType       *string    `json:"reference_type" gorm:"column:reference_type"`
	ReferenceInternalID *int64     `json:"-" gorm:"column:reference_internal_id"`
	ReferencePublicID   *uuid.UUID `json:"reference_public_id" gorm:"column:reference_public_id;type:uuid;index"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type NotificationOutput struct {
	PublicID          uuid.UUID  `json:"public_id" `
	Subject           string     `json:"subject"`
	Message           string     `json:"message"`
	ReferenceType     *string    `json:"reference_type"`
	ReferencePublicID *uuid.UUID `json:"reference_public_id"`
	CreatedAt         time.Time  `json:"created_at"`
}
type NotificationEvent struct {
	RecipientID       int64
	Subject           string
	Message           string
	Type              string
	ReferenceType     string
	ReferenceInternal int64
	ReferencePublic   uuid.UUID
}
