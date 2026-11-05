package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type ResearchGroupRepository interface {
	Create(group *models.ResearchGroup) error
	FetchByPublicID(publicID string) (*models.ResearchGroup, error)
	FetchByJoinCode(code string) (*models.ResearchGroup, error)
	AddMember(member *models.ResearchGroupMember) error
	IsMember(groupID int64, userID int64) (bool, error)
	FetchMembers(groupID int64) ([]models.ResearchGroupMember, error)
	IsSupervisor(groupID int64, userID int64) (bool, error)
	FetchBySlug(slug string) (*models.ResearchGroup, error)
	FetchAll(filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error)
	FetchMyGroups(userID int64, filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error)
	FetchMyGroupByPublicID(userID int64, publicID string) (*models.ResearchGroup, error)
	//Public
	ListPublic(filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error)
	FetchPublicBySlug(slug string) (*models.ResearchGroup, error)
	FetchPublicByPublicID(publicID string) (*models.ResearchGroup, error)
}

type ResearchGroupRepositorys struct {
}

func NewResearchGroupRepository() ResearchGroupRepository {
	return &ResearchGroupRepositorys{}
}
func (r *ResearchGroupRepositorys) Create(group *models.ResearchGroup) error {
	if group.PublicID == uuid.Nil {
		group.PublicID = uuid.New()
	}
	return config.DB.Create(group).Error
}
func (r *ResearchGroupRepositorys) FetchByPublicID(publicID string) (*models.ResearchGroup, error) {
	var group models.ResearchGroup
	err := config.DB.Preload("Leader").Preload("Members.User").Where("public_id = ?", publicID).First(&group).Error
	return &group, err
}

func (r *ResearchGroupRepositorys) FetchByJoinCode(code string) (*models.ResearchGroup, error) {
	var group models.ResearchGroup
	err := config.DB.Where("join_code = ?", code).First(&group).Error
	return &group, err
}

func (r *ResearchGroupRepositorys) AddMember(member *models.ResearchGroupMember) error {
	if member.PublicID == uuid.Nil {
		member.PublicID = uuid.New()
	}
	return config.DB.Create(member).Error
}

func (r *ResearchGroupRepositorys) IsMember(groupID int64, userID int64) (bool, error) {
	var count int64
	err := config.DB.Model(&models.ResearchGroupMember{}).Where("group_internal_id = ? AND user_internal_id = ? AND left_at IS NULL", groupID, userID).Count(&count).Error
	return count > 0, err
}

func (r *ResearchGroupRepositorys) FetchMembers(groupID int64) ([]models.ResearchGroupMember, error) {
	var members []models.ResearchGroupMember
	err := config.DB.Preload("User").Where("group_internal_id = ?", groupID).Order("joined_at ASC").Find(&members).Error
	return members, err
}

// IsSupervisor reports whether the user is the BRIDA reviewer (pembimbing) of the group's research.
func (r *ResearchGroupRepositorys) IsSupervisor(groupID int64, userID int64) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Research{}).
		Where("group_internal_id = ? AND reviewer_internal_id = ? AND deleted_at IS NULL", groupID, userID).Count(&count).Error
	return count > 0, err
}
func (r *ResearchGroupRepositorys) FetchBySlug(slug string) (*models.ResearchGroup, error) {
	var group models.ResearchGroup
	err := config.DB.Preload("Leader").Preload("Members.User").Where("slug = ?", slug).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *ResearchGroupRepositorys) FetchAll(filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error) {
	var groups []models.ResearchGroup
	var total int64
	db := config.DB.Model(&models.ResearchGroup{})
	if filter != "" {
		pattern := "%" + filter + "%"
		db = db.Where("name ILIKE ?", pattern)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	switch sort {
	case "oldest":
		db = db.Order("created_at ASC")
	case "name":
		db = db.Order("name ASC")
	default:
		db = db.Order("created_at DESC")
	}
	err := db.Preload("Leader").Preload("Members").Offset(offset).Limit(limit).Find(&groups).Error
	return groups, total, err
}
func (r *ResearchGroupRepositorys) ListPublic(filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error) {
	var groups []models.ResearchGroup
	var total int64

	db := config.DB.Model(&models.ResearchGroup{}).
		Joins("JOIN researches ON researches.research_id = research_groups.research_internal_id").
		Where("researches.status = ?", models.SubmissionOnGoing)

	if filter != "" {
		pattern := "%" + filter + "%"
		db = db.Where("research_groups.name ILIKE ?", pattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "oldest":
		db = db.Order("research_groups.created_at ASC")
	case "name":
		db = db.Order("research_groups.name ASC")
	default:
		db = db.Order("research_groups.created_at DESC")
	}

	err := db.
		Preload("Leader").
		Preload("Research").
		Offset(offset).
		Limit(limit).
		Find(&groups).Error

	return groups, total, err
}

func (r *ResearchGroupRepositorys) FetchMyGroups(userID int64, filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error) {
	var groups []models.ResearchGroup
	var total int64

	db := config.DB.Model(&models.ResearchGroup{}).
		Select("research_groups.*").
		Joins("JOIN research_group_members ON research_group_members.group_internal_id = research_groups.group_id").
		Where("research_group_members.user_internal_id = ?", userID)

	if filter != "" {
		pattern := "%" + filter + "%"
		db = db.Where("research_groups.name ILIKE ?", pattern)
	}

	switch sort {
	case "oldest":
		db = db.Order("research_groups.created_at ASC")
	case "name":
		db = db.Order("research_groups.name ASC")
	default:
		db = db.Order("research_groups.created_at DESC")
	}

	err := db.
		Preload("Leader").
		Preload("Research").
		Preload("Members").
		Offset(offset).
		Limit(limit).
		Find(&groups).Error

	return groups, total, err
}

func (r *ResearchGroupRepositorys) FetchMyGroupByPublicID(userID int64, publicID string) (*models.ResearchGroup, error) {
	var group models.ResearchGroup
	err := config.DB.
		Select("research_groups.*").
		Joins("JOIN research_group_members ON research_group_members.group_internal_id = research_groups.group_id").
		Where("research_groups.public_id = ? AND research_group_members.user_internal_id = ?", publicID, userID).
		Preload("Leader").
		Preload("Research").
		Preload("Members.User").
		First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *ResearchGroupRepositorys) FetchPublicBySlug(slug string) (*models.ResearchGroup, error) {
	var group models.ResearchGroup
	err := config.DB.
		Joins("JOIN researches ON researches.research_id = research_groups.research_internal_id").
		Where("research_groups.slug = ? AND researches.status = ?", slug, models.SubmissionOnGoing).
		Preload("Leader").
		Preload("Research").
		First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *ResearchGroupRepositorys) FetchPublicByPublicID(publicID string) (*models.ResearchGroup, error) {
	var group models.ResearchGroup
	err := config.DB.
		Joins("JOIN researches ON researches.research_id = research_groups.research_internal_id").
		Where("research_groups.public_id = ? AND researches.status = ?", publicID, models.SubmissionOnGoing).
		Preload("Leader").
		Preload("Research").
		First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}
