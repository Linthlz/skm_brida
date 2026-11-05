package repositories

import (
	"errors"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type NewsRepository interface {
	Create(news *models.News) error
	FetchByPublicID(publicID string) (*models.News, error)
	FetchBySlug(slug string) (*models.News, error)
	SlugExists(slug string) (bool, error)
	List(filter, categorySlug, sort string, offset, limit int, publishedOnly bool) ([]models.News, int64, error)
	Update(news *models.News) error
	SoftDelete(publicID string) error
}

type NewsRepositorys struct{}

func NewNewsRepository() NewsRepository {
	return &NewsRepositorys{}
}

func (r *NewsRepositorys) Create(news *models.News) error {
	return config.DB.Create(news).Error
}

func (r *NewsRepositorys) FetchByPublicID(publicID string) (*models.News, error) {
	var news models.News
	err := config.DB.Preload("Author").Preload("Category").Where("public_id = ?", publicID).First(&news).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *NewsRepositorys) FetchBySlug(slug string) (*models.News, error) {
	var news models.News
	err := config.DB.Preload("Author").Preload("Category").Where("slug = ?", slug).First(&news).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *NewsRepositorys) SlugExists(slug string) (bool, error) {
	var count int64
	err := config.DB.Unscoped().Model(&models.News{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *NewsRepositorys) List(filter, categorySlug, sort string, offset, limit int, publishedOnly bool) ([]models.News, int64, error) {
	var list []models.News
	var total int64
	db := config.DB.Model(&models.News{})
	if publishedOnly {
		db = db.Where("status = ?", models.NewsStatusPublished)
	}
	if filter != "" {
		pattern := "%" + filter + "%"
		db = db.Where("title ILIKE ?", pattern)
	}
	if categorySlug != "" {
		db = db.Joins("JOIN news_categories ON news_categories.category_id = news.category_internal_id").
			Where("news_categories.slug = ?", categorySlug)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	switch sort {
	case "oldest":
		db = db.Order("published_at ASC NULLS LAST")
	case "title":
		db = db.Order("title ASC")
	default:
		db = db.Order("published_at DESC NULLS LAST, created_at DESC")
	}
	err := db.Preload("Author").Preload("Category").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *NewsRepositorys) Update(news *models.News) error {
	return config.DB.Model(&models.News{}).Where("news_id = ?", news.InternalID).Updates(map[string]interface{}{
		"title":                news.Title,
		"category_internal_id": news.CategoryInternalID,
		"content":              news.Content,
		"cover_image":          news.CoverImage,
		"source":               news.Source,
		"status":               news.Status,
		"published_at":         news.PublishedAt,
	}).Error
}

func (r *NewsRepositorys) SoftDelete(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.News{}).Error
}
