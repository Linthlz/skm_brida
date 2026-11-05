package services

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"gorm.io/gorm"
)

const maxCategoryNameLength = 100

type NewsCategoryService interface {
	Create(input models.NewsCategory) (*models.NewsCategory, error)
	List(activeOnly bool) ([]models.NewsCategory, error)
	Update(publicID string, input models.NewsCategory) (*models.NewsCategory, error)
	Delete(publicID string) error
}

type NewsCategoryServices struct {
	NCRepo repositories.NewsCategoryRepository
	db     *gorm.DB
}

func NewNewsCategoryService(NCRepo repositories.NewsCategoryRepository, db *gorm.DB) NewsCategoryService {
	return &NewsCategoryServices{NCRepo: NCRepo, db: db}
}

func validateCategoryName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", utils.ErrValidation("name is required")
	}
	if utf8.RuneCountInString(name) > maxCategoryNameLength {
		return "", utils.ErrValidation(fmt.Sprintf("name must be at most %d characters", maxCategoryNameLength))
	}
	return name, nil
}

func (s *NewsCategoryServices) Create(input models.NewsCategory) (*models.NewsCategory, error) {
	name, err := validateCategoryName(input.Name)
	if err != nil {
		return nil, err
	}
	input.Name = name
	slug := utils.SlugifyResearch(input.Name)
	exists, err := s.NCRepo.SlugExists(slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrConflict("category already exists")
	}
	input.PublicID = uuid.New()
	input.Slug = slug
	if input.DisplayOrder == 0 {
		input.DisplayOrder = 999
	}
	if err := s.NCRepo.Create(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *NewsCategoryServices) List(activeOnly bool) ([]models.NewsCategory, error) {
	return s.NCRepo.List(activeOnly)
}

func (s *NewsCategoryServices) Update(publicID string, input models.NewsCategory) (*models.NewsCategory, error) {
	cat, err := s.NCRepo.FetchByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, utils.ErrNotFound("category not found")
	}
	name, err := validateCategoryName(input.Name)
	if err != nil {
		return nil, err
	}
	newSlug := utils.SlugifyResearch(name)
	if newSlug != cat.Slug {
		exists, err := s.NCRepo.SlugExists(newSlug)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, utils.ErrConflict("category with this name already exists")
		}
	}
	cat.Name = name
	cat.Slug = newSlug
	cat.DisplayOrder = input.DisplayOrder
	cat.IsActive = input.IsActive
	if err := s.NCRepo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *NewsCategoryServices) Delete(publicID string) error {
	cat, err := s.NCRepo.FetchByPublicID(publicID)
	if err != nil {
		return err
	}
	if cat == nil {
		return utils.ErrNotFound("category not found")
	}
	return s.NCRepo.SoftDelete(publicID)
}
