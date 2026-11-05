package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	MonevStatusDraft  = "draft"
	MonevStatusOpen   = "open"
	MonevStatusClosed = "closed"

	OpdEntryStatusPending   = "pending"
	OpdEntryStatusSubmitted = "submitted"
	OpdEntryStatusVerified  = "verified"
)

type Opd struct {
	InternalID int64           `json:"internal_id" gorm:"primaryKey;column:opd_id"`
	PublicID   uuid.UUID       `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	Name       string          `json:"name" gorm:"not null"`
	Code       *string         `json:"code,omitempty" gorm:"column:code;uniqueIndex"`
	Address    string          `json:"address"`
	Phone      string          `json:"phone"`
	IsActive   bool            `json:"is_active" gorm:"column:is_active;default:true"`
	Users      []User          `json:"-" gorm:"foreignKey:OpdInternalID;references:InternalID"`
	Entries    []MonevOpdEntry `json:"-" gorm:"foreignKey:OpdInternalID;references:InternalID"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `json:"-" gorm:"index"`
}

type MonevBatch struct {
	InternalID          int64           `json:"internal_id" gorm:"primaryKey;column:batch_id"`
	PublicID            uuid.UUID       `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	Year                int             `json:"year" gorm:"not null;uniqueIndex"`
	Title               string          `json:"title" gorm:"not null"`
	Description         string          `json:"description"`
	Status              string          `json:"status" gorm:"not null;default:'draft'"`
	CreatedByInternalID int64           `json:"-" gorm:"column:created_by_internal_id;not null"`
	CreatedBy           *User           `json:"created_by,omitempty" gorm:"foreignKey:CreatedByInternalID;references:InternalID"`
	Entries             []MonevOpdEntry `json:"entries,omitempty" gorm:"foreignKey:BatchInternalID;references:InternalID"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
	DeletedAt           gorm.DeletedAt  `json:"-" gorm:"index"`
}

type MonevOpdEntry struct {
	InternalID         int64            `json:"internal_id" gorm:"primaryKey;column:entry_id"`
	PublicID           uuid.UUID        `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	BatchInternalID    int64            `json:"-" gorm:"column:batch_internal_id;not null;index"`
	Batch              *MonevBatch      `json:"batch,omitempty" gorm:"foreignKey:BatchInternalID;references:InternalID"`
	OpdInternalID      int64            `json:"-" gorm:"column:opd_internal_id;not null;index"`
	Opd                *Opd             `json:"opd,omitempty" gorm:"foreignKey:OpdInternalID;references:InternalID"`
	RepresentativeName string           `json:"representative_name"`
	RepresentativeNip  string           `json:"representative_nip"`
	Status             string           `json:"status" gorm:"not null;default:'pending'"`
	FilledByInternalID *int64           `json:"-" gorm:"column:filled_by_internal_id"`
	FilledBy           *User            `json:"filled_by,omitempty" gorm:"foreignKey:FilledByInternalID;references:InternalID"`
	SubmittedAt        *time.Time       `json:"submitted_at,omitempty" gorm:"column:submitted_at"`
	VerifiedAt         *time.Time       `json:"verified_at,omitempty" gorm:"column:verified_at"`
	Kajians            []KajianResearch `json:"kajians,omitempty" gorm:"foreignKey:OpdEntryInternalID;references:InternalID"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	DeletedAt          gorm.DeletedAt   `json:"-" gorm:"index"`
}

type KajianResearch struct {
	InternalID         int64           `json:"internal_id" gorm:"primaryKey;column:kajian_id"`
	PublicID           uuid.UUID       `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	OpdEntryInternalID int64           `json:"-" gorm:"column:opd_entry_internal_id;not null;index"`
	OpdEntry           *MonevOpdEntry  `json:"-" gorm:"foreignKey:OpdEntryInternalID;references:InternalID"`
	Title              string          `json:"title" gorm:"not null"`
	Description        string          `json:"description"`
	DisplayOrder       int             `json:"display_order" gorm:"column:display_order;default:0"`
	Points             []ResearchPoint `json:"points,omitempty" gorm:"foreignKey:KajianInternalID;references:InternalID"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `json:"-" gorm:"index"`
}

type ResearchPoint struct {
	InternalID                    int64          `json:"internal_id" gorm:"primaryKey;column:point_id"`
	PublicID                      uuid.UUID      `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	KajianInternalID              int64          `json:"-" gorm:"column:kajian_internal_id;not null;index"`
	Title                         string         `json:"title" gorm:"not null"`
	MonitoringFilePath            *string        `json:"monitoring_file_path,omitempty" gorm:"column:monitoring_file_path"`
	HasMonitoringResult           bool           `json:"has_monitoring_result" gorm:"column:has_monitoring_result;default:false"`
	MonitoringResultDescription   string         `json:"monitoring_result_description" gorm:"column:monitoring_result_description"`
	KendalaPelaksanaanRekomendasi string         `json:"kendala_pelaksanaan_rekomendasi" gorm:"column:kendala_pelaksanaan_rekomendasi"`
	Benefits                      string         `json:"benefits"`
	DisplayOrder                  int            `json:"display_order" gorm:"column:display_order;default:0"`
	CreatedAt                     time.Time      `json:"created_at"`
	UpdatedAt                     time.Time      `json:"updated_at"`
	DeletedAt                     gorm.DeletedAt `json:"-" gorm:"index"`
}
