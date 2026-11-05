package repositories

import (
	"errors"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type NewsCategoryRepository interface {
	Create(cat *models.NewsCategory) error
	FetchByPublicID(publicID string) (*models.NewsCategory, error)
	FetchBySlug(slug string) (*models.NewsCategory, error)
	List(activeOnly bool) ([]models.NewsCategory, error)
	SlugExists(slug string) (bool, error)
	Update(cat *models.NewsCategory) error
	SoftDelete(publicID string) error
}

type NewsCategoryRepositorys struct{}

func NewNewsCategoryRepository() NewsCategoryRepository {
	return &NewsCategoryRepositorys{}
}

func (r *NewsCategoryRepositorys) Create(cat *models.NewsCategory) error {
	return config.DB.Create(cat).Error
}

func (r *NewsCategoryRepositorys) FetchByPublicID(publicID string) (*models.NewsCategory, error) {
	var cat models.NewsCategory
	err := config.DB.Where("public_id = ?", publicID).First(&cat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *NewsCategoryRepositorys) FetchBySlug(slug string) (*models.NewsCategory, error) {
	var cat models.NewsCategory
	err := config.DB.Where("slug = ?", slug).First(&cat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *NewsCategoryRepositorys) List(activeOnly bool) ([]models.NewsCategory, error) {
	var cats []models.NewsCategory
	db := config.DB.Model(&models.NewsCategory{})
	if activeOnly {
		db = db.Where("is_active = ?", true)
	}
	err := db.Order("display_order ASC, name ASC").Find(&cats).Error
	return cats, err
}

func (r *NewsCategoryRepositorys) SlugExists(slug string) (bool, error) {
	var count int64
	err := config.DB.Unscoped().Model(&models.NewsCategory{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *NewsCategoryRepositorys) Update(cat *models.NewsCategory) error {
	return config.DB.Model(&models.NewsCategory{}).Where("category_id = ?", cat.InternalID).Updates(map[string]interface{}{
		"name":          cat.Name,
		"slug":          cat.Slug,
		"display_order": cat.DisplayOrder,
		"is_active":     cat.IsActive,
	}).Error
}

func (r *NewsCategoryRepositorys) SoftDelete(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.NewsCategory{}).Error
}
