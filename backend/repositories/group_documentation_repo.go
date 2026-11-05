package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type GroupDocumentationRepository interface {
	Create(doc *models.GroupDocumentation, images []models.DocumentationImage) error
	FetchByPublicID(publicID string) (*models.GroupDocumentation, error)
	FetchByGroup(groupInternalID int64, offset, limit int) ([]models.GroupDocumentation, int64, error)
	Update(doc *models.GroupDocumentation) error
	SoftDelete(publicID string) error
	AddImage(image *models.DocumentationImage) error
	DeleteImage(imagePublicID string) error
}

type GroupDocumentationRepositorys struct {
}

func NewGroupDocumentationRepository() GroupDocumentationRepository {
	return &GroupDocumentationRepositorys{}
}
func (r *GroupDocumentationRepositorys) Create(doc *models.GroupDocumentation, images []models.DocumentationImage) error {
	if doc.PublicID == uuid.Nil {
		doc.PublicID = uuid.New()
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(doc).Error; err != nil {
			return err
		}
		for i := range images {
			images[i].DocumentationInternalID = doc.InternalID
			if images[i].PublicID == uuid.Nil {
				images[i].PublicID = uuid.New()
			}
			images[i].DisplayOrder = i
		}
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GroupDocumentationRepositorys) FetchByPublicID(publicID string) (*models.GroupDocumentation, error) {
	var doc models.GroupDocumentation
	err := config.DB.Preload("Author").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Where("public_id = ?", publicID).First(&doc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *GroupDocumentationRepositorys) FetchByGroup(groupInternalID int64, offset, limit int) ([]models.GroupDocumentation, int64, error) {
	var docs []models.GroupDocumentation
	var total int64
	db := config.DB.Model(&models.GroupDocumentation{}).Where("group_internal_id = ?", groupInternalID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Preload("Author").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Order("event_date DESC NULLS LAST, created_at DESC").Offset(offset).Limit(limit).Find(&docs).Error
	return docs, total, err
}

func (r *GroupDocumentationRepositorys) Update(doc *models.GroupDocumentation) error {
	return config.DB.Model(&models.GroupDocumentation{}).Where("documentation_id = ?", doc.InternalID).Updates(map[string]interface{}{
		"title":      doc.Title,
		"narrative":  doc.Narrative,
		"location":   doc.Location,
		"video_url":  doc.VideoURL,
		"event_date": doc.EventDate,
	}).Error
}

func (r *GroupDocumentationRepositorys) SoftDelete(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.GroupDocumentation{}).Error
}

func (r *GroupDocumentationRepositorys) AddImage(image *models.DocumentationImage) error {
	if image.PublicID == uuid.Nil {
		image.PublicID = uuid.New()
	}
	return config.DB.Create(image).Error
}

func (r *GroupDocumentationRepositorys) DeleteImage(imagePublicID string) error {
	return config.DB.Where("public_id = ?", imagePublicID).Delete(&models.DocumentationImage{}).Error
}
