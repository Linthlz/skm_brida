package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status that can be Used Approved|Revision|Rejected
type SubmissionComment struct {
	InternalID           int64     `json:"internal_id" gorm:"primaryKey;column:comment_id"`
	PublicID             uuid.UUID `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	SubmissionInternalID int64     `json:"-" gorm:"column:submission_internal_id;not null;index"`
	Submission           Research  `json:"-" gorm:"foreignKey:SubmissionInternalID;references:InternalID"`

	AuthorInternalID int64 `json:"-" gorm:"column:author_internal_id;not null;index"`
	Author           User  `json:"author,omitempty" gorm:"foreignKey:AuthorInternalID;references:InternalID"`

	Message string `json:"message" gorm:"type:text;not null"`
	// Decision string `json:"decision,omitempty" gorm:"column:decision;type:varchar(50)"`
	Decision *SubmissionDecision `json:"decision,omitempty" gorm:"column:decision;type:varchar(50);check:decision IN ('approved','rejected','revision')"`
	Round    *int                `json:"round,omitempty" gorm:"column:round"`

	ParentCommentID *int64 `json:"-" gorm:"column:parent_comment_id"`
	// Parent          *SubmissionComment `json:"parent,omitempty" gorm:"foreignKey:ParentCommentI"`
	Parent    *SubmissionComment `json:"parent,omitempty" gorm:"foreignKey:ParentCommentID;references:InternalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `json:"-" gorm:"index"`
}
type SubmissionDecision string

const (
	DecisionApproved SubmissionDecision = "approved"
	DecisionRejected SubmissionDecision = "rejected"
	DecisionRevision SubmissionDecision = "revision"
)

type CommentAuthor struct {
	Name *string `json:"name"`
}
type Authors struct {
	PublicID uuid.UUID `json:"public_id"`
	Name     *string   `json:"name"`
	Email    string    `json:"email"`
	// Phone_Number        string    `json:"phone_number"`
	Institution         string `json:"institution"`
	Institution_Address string `json:"institution_address"`
	Position            string `json:"position"`
	// Place_Of_Origin     string         `json:"place_of_origin"`
	Role       string    `json:"role"`
	Verified   bool      `json:"verified"`
	VerifiedAt time.Time `json:"verified_at"`
}
type SubmissionCommentOutput struct {
	PublicID         uuid.UUID `json:"public_id"`
	AuthorPublicID   uuid.UUID `json:"author_public_id"`
	AuthorName       string    `json:"author_name"`
	AuthorPosition   string    `json:"author_position"`
	Content          string    `json:"content"`
	Decision         *string   `json:"decision,omitempty"`
	Round            *int      `json:"round,omitempty"`
	ParentCommentPID *string   `json:"parent_comment_public_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
