package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	FundingOpen    = "open"
	FundingClosing = "closing"
	FundingSoon    = "soon"
)

type FundingScheme struct {
	InternalID   int64          `json:"-" gorm:"primaryKey;column:funding_id"`
	PublicID     uuid.UUID      `json:"public_id" gorm:"column:public_id"`
	Name         string         `json:"name"`
	Organizer    string         `json:"organizer"`
	SchemeCode   string         `json:"scheme_code" gorm:"column:scheme_code"`
	Ceiling      int64          `json:"ceiling"`
	Quota        int            `json:"quota"`
	Deadline     time.Time      `json:"deadline" gorm:"type:date"`
	Status       string         `json:"status"`
	TargetFields string         `json:"target_fields" gorm:"column:target_fields"`
	Requirements string         `json:"requirements"`
	Description  string         `json:"description"`
	Website      string         `json:"website"`
	WebsiteName  string         `json:"website_name" gorm:"column:website_name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type SiteSetting struct {
	Key       string    `gorm:"primaryKey;column:setting_key"`
	Value     string    `gorm:"column:setting_value"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
