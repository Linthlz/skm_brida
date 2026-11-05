package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	InternalID   int64     `json:"internal_id" db:"user_Id" gorm:"primaryKey"`
	PublicID     uuid.UUID `json:"public_id" db:"public_id"`
	Name         *string   `json:"name" db:"name"`
	Email        string    `json:"-" db:"email" gorm:"unique"`
	Password     string    `json:"-" db:"password" gorm:"column:password"`
	Phone_Number string    `json:"-" db:"phone_number" `
	Institution  string    `json:"institution" db:"institution"`
	Position     string    `json:"position" db:"position"`
	// Place_Of_Origin     string    `json:"place_of_origin" db:"place_of_origin"`
	Institution_Address string     `json:"-" db:"institution_address"`
	Role                string     `json:"role" db:"role"`
	Verified            bool       `json:"-" gorm:"default:false"`
	VerificationToken   *string    `json:"-" gorm:"type:uuid;index"`
	VerificationExpiry  *time.Time `json:"-" gorm:"column:verification_expires_at"`
	// VerificationExpiry
	VerifiedAt    *time.Time `json:"-" db:"verified_at"`
	AuthoredNews  []News     `json:"-" gorm:"foreignKey:AuthorInternalID;references:InternalID"`
	OpdInternalID *int64     `json:"-" gorm:"column:opd_internal_id;index"`
	Opd           *Opd       `json:"opd,omitempty" gorm:"foreignKey:OpdInternalID;references:InternalID"`
	CreatedAt     time.Time  `json:"-" db:"created_at"`
	UpdatedAt     time.Time  `json:"-"`
	DeletedAt     time.Time  `json:"-" gorm:"index"`
}

type UserResponse struct {
	PublicID uuid.UUID `json:"public_id"`
	Name     *string   `json:"name"`
	Email    string    `json:"email"`
	// Phone_Number        string    `json:"phone_number"`
	Institution         string `json:"institution"`
	Institution_Address string `json:"institution_address"`
	Position            string `json:"position"`
	// Place_Of_Origin     string         `json:"place_of_origin"`
	Role                   string                `json:"role"`
	LedGroups              []ResearchGroup       `json:"-" gorm:"foreignKey:LeaderInternalID;references:InternalID"`
	GroupMembers           []ResearchGroupMember `json:"-" gorm:"foreignKey:UserInternalID;references:InternalID"`
	AuthoredDocumentations []GroupDocumentation  `json:"-" gorm:"foreignKey:AuthorInternalID;references:InternalID"`
	Verified               bool                  `json:"verified"`
	VerifiedAt             time.Time             `json:"verified_at"`
	CreatedAt              time.Time             `json:"created_at" `
	UpdatedAt              time.Time             `json:"updated_at"`
	DeletedAt              gorm.DeletedAt        `json:"-" `
}

type UpdateUser struct {
	Name *string `json:"name"`
	// Phone_Number        *string `json:"phone_number"`
	Institution         *string `json:"institution"`
	Institution_Address *string `json:"institution_address"`
	Position            *string `json:"position"`
	// Place_Of_Origin     string         `json:"place_of_origin"`
}

// AdminUpdateUser holds the fields an administrator may change; nil means "leave as is".
type AdminUpdateUser struct {
	Name                *string `json:"name"`
	Institution         *string `json:"institution"`
	Institution_Address *string `json:"institution_address"`
	Position            *string `json:"position"`
	Role                *string `json:"role"`
	Verified            *bool   `json:"verified"`
	Password            *string `json:"password"`
}

type CreateUser struct {
	Name *string `json:"name"`
	// Phone_Number        *string `json:"phone_number"`
	Institution         *string `json:"institution"`
	Institution_Address *string `json:"institution_address"`
	Position            *string `json:"position"`
	Role                string  `json:"role"`
	Email               string  `json:"email"`
	Password            string  `json:"password"`
	// Place_Of_Origin     string         `json:"place_of_origin"`
}
