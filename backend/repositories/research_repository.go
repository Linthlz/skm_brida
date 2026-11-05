package repositories

import (
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
)

type ResearchRepository interface {
	//Basic Funcs
	Submit(research *models.Research) error
	SubmissionUpdate(publicID string, updates map[string]interface{}) error
	Delete(id uint) error
	//Fetch
	ResearchSlug(slug string) (*models.Research, error)
	FetchByPublicID(PubllicID string) (*models.Research, error)
	FetchByID(id uint) (*models.Research, error)
	Researches(filter, sort string, offset, limit int) ([]models.Research, int64, error)
	UserResearches(userInternalID uint, filter, sort string, offset, limit int) ([]models.Research, int64, error)
}

type ResearchRepositorys struct{}

func NewResearchRepository() ResearchRepository {
	return &ResearchRepositorys{}
}

// Basic Func
func (r *ResearchRepositorys) Submit(research *models.Research) error {
	return config.DB.Create(research).Error
}

func (r *ResearchRepositorys) SubmissionUpdate(publicID string, updates map[string]interface{}) error {
	return config.DB.Model(&models.Research{}).Where("public_id=?", publicID).Where("status = ?", "revision").Updates(updates).Error
}

func (r *ResearchRepositorys) Delete(id uint) error {
	return config.DB.Delete(&models.Research{}).Error
}

// Fetch
func (r *ResearchRepositorys) FetchByPublicID(PubllicID string) (*models.Research, error) {
	var research models.Research
	err := config.DB.Where("public_id = ?", PubllicID).First(&research).Error
	return &research, err
}
func (r *ResearchRepositorys) FetchByID(id uint) (*models.Research, error) {
	var research models.Research
	err := config.DB.Where("id = ?", id).First(&research).Error
	return &research, err
}

// For Reviewer/Public
func (r *ResearchRepositorys) Researches(filter, sort string, offset, limit int) ([]models.Research, int64, error) {
	var research []models.Research
	var total int64

	db := config.DB.Model(&models.Research{})
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("title Ilike ?", filterPattern)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	switch sort {
	case "oldest":
		db = db.Order("created_at ASC")
	case "subject":
		db = db.Order("subject ASC")
	default:
		db = db.Order("created_at DESC")
	}
	err := db.Preload("Uploader").Preload("Reviewer").Limit(limit).Offset(offset).Find(&research).Error
	return research, total, err
}

// For personal user
func (r *ResearchRepositorys) UserResearches(userInternalID uint, filter, sort string, offset, limit int) ([]models.Research, int64, error) {
	var research []models.Research
	var total int64

	db := config.DB.Model(&models.Research{}).Where("uploader_internal_id = ?", userInternalID)
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("title Ilike ?", filterPattern)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// if sort != "" {
	// 	if sort == "-id" {
	// 		sort = "-internal_id"
	// 	} else if sort == "id" {
	// 		sort = "internal_id"
	// 	}
	// 	if strings.HasPrefix(sort, "") {
	// 		sort = strings.TrimPrefix(sort, "-") + " DESC"
	// 	} else {
	// 		sort += " ASC"
	// 	}
	// 	db = db.Order(sort)
	// }
	switch sort {
	case "oldest":
		db = db.Order("created_at ASC")
	case "subject":
		db = db.Order("subject ASC")
	default:
		db = db.Order("created_at DESC")
	}
	err := db.Preload("Uploader").Limit(limit).Offset(offset).Find(&research).Error
	return research, total, err
}

func (r *ResearchRepositorys) ResearchSlug(slug string) (*models.Research, error) {
	var submission models.Research
	err := config.DB.Preload("Uploader").Where("slug = ?", slug).First(&submission).Error
	return &submission, err
}

func (r *ResearchRepositorys) Revision(publicID string, updates map[string]interface{}) error {
	return config.DB.Model(&models.Research{}).Where("public_id=?", publicID).Updates(updates).Error
}
