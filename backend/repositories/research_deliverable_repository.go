package repositories

import (
	"errors"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type ResearchDeliverableRepository interface {
	Create(d *models.ResearchDeliverable) error
	FetchByPublicID(publicID string) (*models.ResearchDeliverable, error)
	FetchByResearch(researchInternalID int64) ([]models.ResearchDeliverable, error)
	Update(d *models.ResearchDeliverable) error
	SoftDelete(publicID string) error
	CountByStatus(researchInternalID int64) (total int64, done int64, err error)
	NextDisplayOrder(researchInternalID int64) (int, error)
}

type ResearchDeliverableRepositorys struct{}

func NewResearchDeliverableRepository() ResearchDeliverableRepository {
	return &ResearchDeliverableRepositorys{}
}

func (r *ResearchDeliverableRepositorys) Create(d *models.ResearchDeliverable) error {
	return config.DB.Create(d).Error
}

func (r *ResearchDeliverableRepositorys) FetchByPublicID(publicID string) (*models.ResearchDeliverable, error) {
	var d models.ResearchDeliverable
	err := config.DB.Preload("CompletedBy").Preload("CreatedBy").Where("public_id = ?", publicID).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *ResearchDeliverableRepositorys) FetchByResearch(researchInternalID int64) ([]models.ResearchDeliverable, error) {
	var list []models.ResearchDeliverable
	err := config.DB.Preload("CompletedBy").Where("research_internal_id = ?", researchInternalID).Order("display_order ASC, created_at ASC").Find(&list).Error
	return list, err
}

func (r *ResearchDeliverableRepositorys) Update(d *models.ResearchDeliverable) error {
	return config.DB.Model(&models.ResearchDeliverable{}).Where("deliverable_id = ?", d.InternalID).Updates(map[string]interface{}{
		"title":                    d.Title,
		"description":              d.Description,
		"proof_type":               d.ProofType,
		"status":                   d.Status,
		"file_path":                d.FilePath,
		"completed_by_internal_id": d.CompletedByInternalID,
		"completed_at":             d.CompletedAt,
		"display_order":            d.DisplayOrder,
	}).Error
}

func (r *ResearchDeliverableRepositorys) SoftDelete(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.ResearchDeliverable{}).Error
}

func (r *ResearchDeliverableRepositorys) CountByStatus(researchInternalID int64) (int64, int64, error) {
	var total, done int64
	if err := config.DB.Model(&models.ResearchDeliverable{}).Where("research_internal_id = ?", researchInternalID).Count(&total).Error; err != nil {
		return 0, 0, err
	}
	if err := config.DB.Model(&models.ResearchDeliverable{}).Where("research_internal_id = ? AND status = ?", researchInternalID, models.DeliverableStatusDone).Count(&done).Error; err != nil {
		return 0, 0, err
	}
	return total, done, nil
}

func (r *ResearchDeliverableRepositorys) NextDisplayOrder(researchInternalID int64) (int, error) {
	var max int
	err := config.DB.Model(&models.ResearchDeliverable{}).
		Where("research_internal_id = ?", researchInternalID).
		Select("COALESCE(MAX(display_order), -1)").
		Scan(&max).Error
	return max + 1, err
}
