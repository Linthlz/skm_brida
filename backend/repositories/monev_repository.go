package repositories

import (
	"errors"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type MonevRepository interface {
	CreateBatch(b *models.MonevBatch) error
	FetchBatchByPublicID(publicID string) (*models.MonevBatch, error)
	ListBatches(filter, sort string, offset, limit int) ([]models.MonevBatch, int64, error)
	UpdateBatch(b *models.MonevBatch) error
	SoftDeleteBatch(publicID string) error

	CreateEntry(e *models.MonevOpdEntry) error
	FetchEntryByPublicID(publicID string) (*models.MonevOpdEntry, error)
	ListEntriesByBatch(batchID int64) ([]models.MonevOpdEntry, error)
	ListEntriesByOpd(opdID int64) ([]models.MonevOpdEntry, error)
	UpdateEntry(e *models.MonevOpdEntry) error
	SoftDeleteEntry(publicID string) error

	CreateKajian(k *models.KajianResearch) error
	FetchKajianByPublicID(publicID string) (*models.KajianResearch, error)
	ListKajiansByEntry(entryID int64) ([]models.KajianResearch, error)
	UpdateKajian(k *models.KajianResearch) error
	SoftDeleteKajian(publicID string) error

	CreatePoint(p *models.ResearchPoint) error
	FetchPointByPublicID(publicID string) (*models.ResearchPoint, error)
	ListPointsByKajian(kajianID int64) ([]models.ResearchPoint, error)
	UpdatePoint(p *models.ResearchPoint) error
	SoftDeletePoint(publicID string) error
	NextPointOrder(kajianID int64) (int, error)
}

type MonevRepositorys struct{}

func NewMonevRepository() MonevRepository {
	return &MonevRepositorys{}
}

func (r *MonevRepositorys) CreateBatch(b *models.MonevBatch) error {
	return config.DB.Create(b).Error
}

func (r *MonevRepositorys) FetchBatchByPublicID(publicID string) (*models.MonevBatch, error) {
	var b models.MonevBatch
	err := config.DB.Preload("CreatedBy").Preload("Entries.Opd").Preload("Entries.Kajians.Points").Where("public_id = ?", publicID).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *MonevRepositorys) ListBatches(filter, sort string, offset, limit int) ([]models.MonevBatch, int64, error) {
	var list []models.MonevBatch
	var total int64
	db := config.DB.Model(&models.MonevBatch{})
	if filter != "" {
		pattern := "%" + filter + "%"
		db = db.Where("title ILIKE ?", pattern)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	switch sort {
	case "oldest":
		db = db.Order("year ASC")
	default:
		db = db.Order("year DESC")
	}
	err := db.Preload("CreatedBy").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *MonevRepositorys) UpdateBatch(b *models.MonevBatch) error {
	return config.DB.Model(&models.MonevBatch{}).Where("batch_id = ?", b.InternalID).Updates(map[string]interface{}{
		"year": b.Year, "title": b.Title, "description": b.Description, "status": b.Status,
	}).Error
}

func (r *MonevRepositorys) SoftDeleteBatch(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.MonevBatch{}).Error
}

func (r *MonevRepositorys) CreateEntry(e *models.MonevOpdEntry) error {
	return config.DB.Create(e).Error
}

func (r *MonevRepositorys) FetchEntryByPublicID(publicID string) (*models.MonevOpdEntry, error) {
	var e models.MonevOpdEntry
	err := config.DB.Preload("Batch").Preload("Opd").Preload("Kajians.Points").Where("public_id = ?", publicID).First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *MonevRepositorys) ListEntriesByBatch(batchID int64) ([]models.MonevOpdEntry, error) {
	var list []models.MonevOpdEntry
	err := config.DB.Preload("Opd").Preload("Kajians").Where("batch_internal_id = ?", batchID).Order("created_at ASC").Find(&list).Error
	return list, err
}

func (r *MonevRepositorys) ListEntriesByOpd(opdID int64) ([]models.MonevOpdEntry, error) {
	var list []models.MonevOpdEntry
	err := config.DB.Preload("Batch").Preload("Opd").Preload("Kajians.Points").Where("opd_internal_id = ?", opdID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *MonevRepositorys) UpdateEntry(e *models.MonevOpdEntry) error {
	return config.DB.Model(&models.MonevOpdEntry{}).Where("entry_id = ?", e.InternalID).Updates(map[string]interface{}{
		"representative_name":   e.RepresentativeName,
		"representative_nip":    e.RepresentativeNip,
		"status":                e.Status,
		"filled_by_internal_id": e.FilledByInternalID,
		"submitted_at":          e.SubmittedAt,
		"verified_at":           e.VerifiedAt,
	}).Error
}

func (r *MonevRepositorys) SoftDeleteEntry(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.MonevOpdEntry{}).Error
}

func (r *MonevRepositorys) CreateKajian(k *models.KajianResearch) error {
	return config.DB.Create(k).Error
}

func (r *MonevRepositorys) FetchKajianByPublicID(publicID string) (*models.KajianResearch, error) {
	var k models.KajianResearch
	err := config.DB.Preload("OpdEntry.Batch").Preload("OpdEntry.Opd").Preload("Points").Where("public_id = ?", publicID).First(&k).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *MonevRepositorys) ListKajiansByEntry(entryID int64) ([]models.KajianResearch, error) {
	var list []models.KajianResearch
	err := config.DB.Preload("Points").Where("opd_entry_internal_id = ?", entryID).Order("display_order ASC, created_at ASC").Find(&list).Error
	return list, err
}

func (r *MonevRepositorys) UpdateKajian(k *models.KajianResearch) error {
	return config.DB.Model(&models.KajianResearch{}).Where("kajian_id = ?", k.InternalID).Updates(map[string]interface{}{
		"title": k.Title, "description": k.Description, "display_order": k.DisplayOrder,
	}).Error
}

func (r *MonevRepositorys) SoftDeleteKajian(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.KajianResearch{}).Error
}

func (r *MonevRepositorys) CreatePoint(p *models.ResearchPoint) error {
	return config.DB.Create(p).Error
}

func (r *MonevRepositorys) FetchPointByPublicID(publicID string) (*models.ResearchPoint, error) {
	var p models.ResearchPoint
	err := config.DB.Where("public_id = ?", publicID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *MonevRepositorys) ListPointsByKajian(kajianID int64) ([]models.ResearchPoint, error) {
	var list []models.ResearchPoint
	err := config.DB.Where("kajian_internal_id = ?", kajianID).Order("display_order ASC, created_at ASC").Find(&list).Error
	return list, err
}

func (r *MonevRepositorys) UpdatePoint(p *models.ResearchPoint) error {
	return config.DB.Model(&models.ResearchPoint{}).Where("point_id = ?", p.InternalID).Updates(map[string]interface{}{
		"title":                           p.Title,
		"monitoring_file_path":            p.MonitoringFilePath,
		"has_monitoring_result":           p.HasMonitoringResult,
		"monitoring_result_description":   p.MonitoringResultDescription,
		"kendala_pelaksanaan_rekomendasi": p.KendalaPelaksanaanRekomendasi,
		"benefits":                        p.Benefits,
		"display_order":                   p.DisplayOrder,
	}).Error
}

func (r *MonevRepositorys) SoftDeletePoint(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.ResearchPoint{}).Error
}

func (r *MonevRepositorys) NextPointOrder(kajianID int64) (int, error) {
	var max int
	err := config.DB.Model(&models.ResearchPoint{}).Where("kajian_internal_id = ?", kajianID).Select("COALESCE(MAX(display_order), -1)").Scan(&max).Error
	return max + 1, err
}
