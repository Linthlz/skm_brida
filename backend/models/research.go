package models

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionStatus string

const (
	SubmissionPending     SubmissionStatus = "pending"
	SubmissionUnderReview SubmissionStatus = "under-review"
	SubmissionOnGoing     SubmissionStatus = "on-going"
	SubmissionApproved    SubmissionStatus = "approved"
	SubmissionRejected    SubmissionStatus = "rejected"
	SubmissionRevision    SubmissionStatus = "revision"
)

type Research struct {
	InternalID      *int64    `json:"internal_id" gorm:"primaryKey;column:research_id"`
	PublicID        uuid.UUID `json:"public_id" db:"public_id"`
	ResearchCode    *string   `json:"research_code" gorm:"column:research_code;uniqueIndex;not null"`
	Title           *string   `json:"title" db:"title"`
	Slug            string    `json:"slug" gorm:"column:slug;uniqueIndex"`
	Purpose         string    `json:"Purpose" db:"Purpose" `
	Category        string    `json:"category" db:"category" `
	Funding_Scheme  string    `json:"funding_scheme" db:"funding_scheme" `
	Fund_Amount     string    `json:"fund_amount" db:"fund_amount" `
	Location        string    `json:"location" db:"location"`
	Address         string    `json:"address" db:"address"`
	RPJMD           string    `json:"rpjmd" db:"rpjmd"`
	Promised_Output string    `json:"promised_output" db:"promised_output"`
	Abstract        string    `json:"abstract,omitempty" gorm:"column:abstract;type:text"`
	Description     string    `json:"description,omitempty" gorm:"column:description;type:text"`
	// Fields Part 2
	Research_Significance                    string `json:"research_significance" db:"research_significance"`
	Research_Time_Range                      string `json:"research_time_range" db:"research_time_range"`
	Partnership_Target_Selection             string `json:"partnership_target_selection" db:"partnership_target_selection"`
	Partner_Name                             string `json:"partner_name" db:"partner_name"`
	Estimated_Number_of_Direct_Beneficiaries string `json:"estimated_number_of_direct_beneficiaries" db:"estimated_number_of_direct_beneficiaries"`
	// Field Part 3 (File Upload)
	File string `json:"research_proposal_file" gorm:"column:research_proposal_file"`
	// Status
	Status             SubmissionStatus `json:"status" gorm:"default:pending"`
	UploaderInternalID *int64           `json:"-" gorm:"column:uploader_internal_id"`
	Uploader           *User            `json:"uploader,omitempty" gorm:"foreignKey:UploaderInternalID;references:InternalID"`
	ReviewerInternalID *int64           `json:"-" gorm:"column:reviewer_internal_id"`
	Reviewer           *User            `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerInternalID;references:InternalID"`
	// UploaderPublicID   uuid.UUID        `json:"-" gorm:"column:uploader_public_id"`
	GroupInternalID *int64                `json:"-" gorm:"column:group_internal_id;index"`
	Group           *ResearchGroup        `json:"group,omitempty" gorm:"foreignKey:GroupInternalID;references:InternalID"`
	Deliverables    []ResearchDeliverable `json:"deliverables,omitempty" gorm:"foreignKey:ResearchInternalID;references:InternalID"`
	Documentations  []GroupDocumentation  `json:"-" gorm:"foreignKey:GroupInternalID;references:InternalID"`
	CreatedAt       time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	DeletedAt       *time.Time            `json:"-" gorm:"index"`
}

type UpdateSubmission struct {
	Title           *string `json:"title"`
	Purpose         *string `json:"purpose" `
	Category        *string `json:"category"  `
	Funding_Scheme  *string `json:"funding_scheme"  `
	Fund_Amount     *string `json:"fund_amount" `
	Location        *string `json:"location" `
	Address         *string `json:"address" `
	RPJMD           *string `json:"rpjmd"`
	Promised_Output *string `json:"promised_output"`
	// Fields Part 2
	Research_Significance                    *string `json:"research_significance"`
	Research_Time_Range                      *string `json:"research_time_range" `
	Partnership_Target_Selection             *string `json:"partnership_target_selection" `
	Partner_Name                             *string `json:"partner_name"`
	Estimated_Number_of_Direct_Beneficiaries *string `json:"estimated_number_of_direct_beneficiaries"`
	// Field Part 3 (File Upload)
	File *string `json:"research_proposal_file" `
}
type SubmissionResponse struct {
	PublicID            uuid.UUID `json:"public_id" `
	ResearchCode        string    `json:"research_code"`
	Title               string    `json:"title"`
	Slug                string    `json:"slug"`
	Category            string    `json:"category"  `
	Location            string    `json:"location" `
	Research_Time_Range string    `json:"research_time_range" `
	File                *string   `json:"research_proposal_file" `
	Status              *string   `json:"status"`
	UploaderPublicID    uuid.UUID `json:"uploader_public_id"`
	UploaderName        string    `json:"uploader_name"`
	FundingScheme       string    `json:"funding_scheme"`
	FundAmount          string    `json:"fund_amount"`
	CreatedAt           time.Time `json:"created_at"`
}
type ReviewerSubmissionResponse struct {
	PublicID            uuid.UUID `json:"public_id" `
	ResearchCode        string    `json:"research_code"`
	Title               string    `json:"title"`
	Slug                string    `json:"slug"`
	Category            string    `json:"category"  `
	Location            string    `json:"location" `
	Research_Time_Range string    `json:"research_time_range" `
	File                *string   `json:"research_proposal_file" `
	Status              *string   `json:"status"`
	ReviewerInternalID  *int64    `json:"reviewer_internal_id "`
	// Reviewer         *UserBrief `json:"reviewer,omitempty"`
	ReviewerName     string    `json:"reviewer_name"`
	ReviewerPosition string    `json:"reviewer_position"`
	UploaderPublicID uuid.UUID `json:"uploader_public_id"`
	UploaderName     string    `json:"uploader_name"`
	UploaderPosition string    `json:"uploader_position"`
}
type Submissions struct {
	PublicID        uuid.UUID `json:"public_id" `
	ResearchCode    string    `json:"research_code"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Purpose         string    `json:"purpose" `
	Category        string    `json:"category"  `
	Funding_Scheme  string    `json:"funding_scheme"  `
	Fund_Amount     string    `json:"fund_amount" `
	Location        string    `json:"location" `
	Address         string    `json:"address" `
	RPJMD           string    `json:"rpjmd"`
	Promised_Output string    `json:"promised_output"`
	// Fields Part 2
	Research_Significance                    string `json:"research_significance"`
	Research_Time_Range                      string `json:"research_time_range" `
	Partnership_Target_Selection             string `json:"partnership_target_selection" `
	Partner_Name                             string `json:"partner_name"`
	Estimated_Number_of_Direct_Beneficiaries string `json:"estimated_number_of_direct_beneficiaries"`
	// Field Part 3 (File Upload)
	File             *string   `json:"research_proposal_file" `
	UploaderPublicID uuid.UUID `json:"uploader_public_id"`
	UploaderName     string    `json:"uploader_name"`
}
