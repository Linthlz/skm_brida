package services

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"gorm.io/gorm"
)

const (
	maxNewsTitleLength  = 255
	maxNewsSourceLength = 255
)

type NewsService interface {
	Create(authorID int64, categoryPublicID string, input models.News) (*models.News, error)
	FetchBySlug(slug string, publishedOnly bool) (*models.News, error)
	FetchByPublicID(publicID string) (*models.News, error)
	List(filter, categorySlug, sort string, offset, limit int, publishedOnly bool) ([]models.News, int64, error)
	Update(authorID int64, publicID string, categoryPublicID string, input models.News) (*models.News, error)
	Delete(authorID int64, publicID string) error
}

type NewsServices struct {
	repo         repositories.NewsRepository
	categoryRepo repositories.NewsCategoryRepository
	db           *gorm.DB
}

func NewNewsService(repo repositories.NewsRepository, categoryRepo repositories.NewsCategoryRepository, db *gorm.DB) NewsService {
	return &NewsServices{repo: repo, categoryRepo: categoryRepo, db: db}
}

func isValidNewsStatus(s string) bool {
	return s == models.NewsStatusDraft || s == models.NewsStatusPublished
}

func generateUniqueNewsSlug(repo repositories.NewsRepository, title string) (string, error) {
	base := utils.SlugifyResearch(title)
	candidate := base
	for i := 1; i <= 5; i++ {
		exists, err := repo.SlugExists(candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return "", utils.ErrConflict("a news item with a similar title already exists")
}

func validateNewsInput(input *models.News) error {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	if input.Title == "" {
		return utils.ErrValidation("title is required")
	}
	if utf8.RuneCountInString(input.Title) > maxNewsTitleLength {
		return utils.ErrValidation(fmt.Sprintf("title must be at most %d characters", maxNewsTitleLength))
	}
	if input.Content == "" {
		return utils.ErrValidation("content is required")
	}
	if input.Source != nil {
		src := strings.TrimSpace(*input.Source)
		if src == "" {
			input.Source = nil
		} else if utf8.RuneCountInString(src) > maxNewsSourceLength {
			return utils.ErrValidation(fmt.Sprintf("source must be at most %d characters", maxNewsSourceLength))
		} else {
			input.Source = &src
		}
	}
	if input.Status != "" && !isValidNewsStatus(input.Status) {
		return utils.ErrValidation("status must be 'draft' or 'published'")
	}
	return nil
}

func (s *NewsServices) resolveCategory(categoryPublicID string) (*models.NewsCategory, error) {
	if _, err := uuid.Parse(categoryPublicID); err != nil {
		return nil, utils.ErrValidation("category_public_id is not a valid id")
	}
	cat, err := s.categoryRepo.FetchByPublicID(categoryPublicID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, utils.ErrValidation("category not found")
	}
	return cat, nil
}

func (s *NewsServices) Create(authorID int64, categoryPublicID string, input models.News) (*models.News, error) {
	if err := validateNewsInput(&input); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.CoverImage) == "" {
		return nil, utils.ErrValidation("cover image is required")
	}
	if strings.TrimSpace(categoryPublicID) == "" {
		return nil, utils.ErrValidation("category is required")
	}
	cat, err := s.resolveCategory(categoryPublicID)
	if err != nil {
		return nil, err
	}
	if input.Status == "" {
		input.Status = models.NewsStatusDraft
	}
	slug, err := generateUniqueNewsSlug(s.repo, input.Title)
	if err != nil {
		return nil, err
	}
	input.PublicID = uuid.New()
	input.Slug = slug
	input.AuthorInternalID = authorID
	input.CategoryInternalID = &cat.InternalID
	if input.Status == models.NewsStatusPublished && input.PublishedAt == nil {
		now := time.Now()
		input.PublishedAt = &now
	}
	if err := s.repo.Create(&input); err != nil {
		return nil, err
	}
	return s.repo.FetchByPublicID(input.PublicID.String())
}

func (s *NewsServices) FetchBySlug(slug string, publishedOnly bool) (*models.News, error) {
	news, err := s.repo.FetchBySlug(slug)
	if err != nil {
		return nil, err
	}
	if news == nil || (publishedOnly && news.Status != models.NewsStatusPublished) {
		return nil, utils.ErrNotFound("news not found")
	}
	return news, nil
}

func (s *NewsServices) FetchByPublicID(publicID string) (*models.News, error) {
	news, err := s.repo.FetchByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, utils.ErrNotFound("news not found")
	}
	return news, nil
}

func (s *NewsServices) List(filter, categorySlug, sort string, offset, limit int, publishedOnly bool) ([]models.News, int64, error) {
	return s.repo.List(filter, categorySlug, sort, offset, limit, publishedOnly)
}

func (s *NewsServices) Update(authorID int64, publicID string, categoryPublicID string, input models.News) (*models.News, error) {
	news, err := s.FetchByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if err := validateNewsInput(&input); err != nil {
		return nil, err
	}
	if strings.TrimSpace(categoryPublicID) != "" {
		cat, err := s.resolveCategory(categoryPublicID)
		if err != nil {
			return nil, err
		}
		news.CategoryInternalID = &cat.InternalID
		news.Category = cat
	}
	if input.Status == "" {
		input.Status = news.Status
	}
	news.Title = input.Title
	news.Content = input.Content
	if strings.TrimSpace(input.CoverImage) != "" {
		news.CoverImage = input.CoverImage
	}
	news.Source = input.Source
	if news.Status != models.NewsStatusPublished && input.Status == models.NewsStatusPublished {
		now := time.Now()
		news.PublishedAt = &now
	}
	news.Status = input.Status
	if err := s.repo.Update(news); err != nil {
		return nil, err
	}
	return s.repo.FetchByPublicID(publicID)
}

func (s *NewsServices) Delete(authorID int64, publicID string) error {
	if _, err := s.FetchByPublicID(publicID); err != nil {
		return err
	}
	return s.repo.SoftDelete(publicID)
}
