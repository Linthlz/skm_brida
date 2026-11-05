package repositories

import (
	"errors"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type OpdRepository interface {
	Create(opd *models.Opd) error
	FetchByPublicID(publicID string) (*models.Opd, error)
	List(activeOnly bool) ([]models.Opd, error)
	Update(opd *models.Opd) error
	SoftDelete(publicID string) error
}

type OpdRepositorys struct{}

func NewOpdRepository() OpdRepository {
	return &OpdRepositorys{}
}

func (r *OpdRepositorys) Create(opd *models.Opd) error {
	return config.DB.Create(opd).Error
}

func (r *OpdRepositorys) FetchByPublicID(publicID string) (*models.Opd, error) {
	var opd models.Opd
	err := config.DB.Where("public_id = ?", publicID).First(&opd).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &opd, nil
}

func (r *OpdRepositorys) List(activeOnly bool) ([]models.Opd, error) {
	var list []models.Opd
	db := config.DB.Model(&models.Opd{})
	if activeOnly {
		db = db.Where("is_active = ?", true)
	}
	err := db.Order("name ASC").Find(&list).Error
	return list, err
}

func (r *OpdRepositorys) Update(opd *models.Opd) error {
	return config.DB.Model(&models.Opd{}).Where("opd_id = ?", opd.InternalID).Updates(map[string]interface{}{
		"name": opd.Name, "code": opd.Code, "address": opd.Address, "phone": opd.Phone, "is_active": opd.IsActive,
	}).Error
}

func (r *OpdRepositorys) SoftDelete(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.Opd{}).Error
}
