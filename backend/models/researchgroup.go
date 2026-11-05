package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResearchGroup struct {
	InternalID         int64     `json:"-" gorm:"primaryKey;column:group_id"`
	PublicID           uuid.UUID `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	ResearchInternalID int64     `json:"-" gorm:"column:research_internal_id;not null;index"`
	Research           *Research `json:"research,omitempty" gorm:"foreignKey:ResearchInternalID;references:InternalID"`
	Name               string    `json:"name" gorm:"not null"`
	Slug               *string   `json:"slug,omitempty" gorm:"column:slug;uniqueIndex"`
	JoinCode           string    `json:"-" gorm:"column:join_code;uniqueIndex;not null"`
	Description        string    `json:"-"`

	LeaderInternalID int64 `json:"-" gorm:"column:leader_internal_id;not null;index"`
	Leader           *User `json:"leader,omitempty" gorm:"foreignKey:LeaderInternalID;references:InternalID"`

	Members []ResearchGroupMember `json:"members,omitempty" gorm:"foreignKey:GroupInternalID;references:InternalID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ResearchGroupMember struct {
	InternalID      int64          `json:"-" gorm:"primaryKey;column:member_id"`
	PublicID        uuid.UUID      `json:"public_id" gorm:"column:public_id;uniqueIndex;not null"`
	GroupInternalID int64          `json:"-" gorm:"column:group_internal_id;not null;index"`
	UserInternalID  int64          `json:"-" gorm:"column:user_internal_id;not null;index"`
	User            *User          `json:"user,omitempty" gorm:"foreignKey:UserInternalID;references:InternalID"`
	Role            string         `json:"role" gorm:"not null;default:'member'"`
	JoinedAt        time.Time      `json:"joined_at"`
	LeftAt          *time.Time     `json:"left_at,omitempty"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}
type PublicGroupOutput struct {
	PublicID    uuid.UUID `json:"public_id"`
	Name        string    `json:"name"`
	Slug        *string   `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`

	Leader struct {
		PublicID uuid.UUID `json:"public_id"`
		Name     string    `json:"name"`
		Position string    `json:"position"`
	} `json:"leader"`

	MemberCount int `json:"member_count"`

	Research *PublicResearchOutput `json:"research,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type PublicResearchOutput struct {
	PublicID              uuid.UUID `json:"public_id"`
	Slug                  *string   `json:"slug,omitempty"`
	Title                 string    `json:"title"`
	Category              string    `json:"category"`
	Location              string    `json:"location"`
	Address               string    `json:"address,omitempty"`
	Research_Time_Range   string    `json:"research_time_range"`
	Status                string    `json:"status"`
	Purpose               string    `json:"purpose,omitempty"`
	Research_Significance string    `json:"research_significance,omitempty"`
	Promised_Output       string    `json:"promised_output,omitempty"`
}
type MyGroupDetailOutput struct {
	PublicID    uuid.UUID `json:"public_id"`
	Name        string    `json:"name"`
	Slug        *string   `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
	JoinCode    string    `json:"join_code"`
	MyRole      string    `json:"my_role"`

	Leader struct {
		PublicID uuid.UUID `json:"public_id"`
		Name     string    `json:"name"`
		Position string    `json:"position"`
	} `json:"leader"`

	Members []struct {
		PublicID uuid.UUID `json:"public_id"`
		Name     string    `json:"name"`
		Position string    `json:"position"`
		Role     string    `json:"role"`
		JoinedAt time.Time `json:"joined_at"`
	} `json:"members"`

	Research  *PublicResearchOutput `json:"research,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
}
