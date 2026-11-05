package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"gorm.io/gorm"
)

type GroupDocumentationService interface {
	Create(authorID int64, groupPublicID string, input models.GroupDocumentation, imagePaths []string, imageCaptions []string) (*models.GroupDocumentation, error)
	FetchByPublicID(publicID string) (*models.GroupDocumentation, error)
	FetchByGroup(groupPublicID string, offset, limit int) ([]models.GroupDocumentation, int64, error)
	AddImage(authorID int64, docPublicID string, filePath string, caption string) (*models.DocumentationImage, error)
	Delete(authorID int64, publicID string) error
	DeleteImage(authorID int64, docPublicID string, imagePublicID string) error
	Update(authorID int64, publicID string, input models.GroupDocumentation) error
}
type GroupDocumentationServices struct {
	GDRepo repositories.GroupDocumentationRepository
	GRepo  repositories.ResearchGroupRepository
	db     *gorm.DB
}

func NewGroupDocumentationService(GDRepo repositories.GroupDocumentationRepository, GRepo repositories.ResearchGroupRepository, db *gorm.DB) GroupDocumentationService {
	return &GroupDocumentationServices{GDRepo: GDRepo, GRepo: GRepo, db: db}
}

// canManage allows group members and global admins to manage documentation.
func (s *GroupDocumentationServices) canManage(groupID, userID int64) (bool, error) {
	ok, err := s.GRepo.IsMember(groupID, userID)
	if err != nil || ok {
		return ok, err
	}
	var role string
	if err := s.db.Table("users").Select("role").Where("internal_id = ?", userID).Scan(&role).Error; err != nil {
		return false, err
	}
	return role == "admin", nil
}

func (s *GroupDocumentationServices) Create(authorID int64, groupPublicID string, input models.GroupDocumentation, imagePaths []string, imageCaptions []string) (*models.GroupDocumentation, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title is required")
	}
	if len(imagePaths) == 0 {
		return nil, errors.New("at least one image is required")
	}
	if len(imagePaths) > utils.MaxImagesPerDocumentation {
		return nil, fmt.Errorf("maximum %d images per documentation", utils.MaxImagesPerDocumentation)
	}
	cleanURL, err := utils.NormalizeVideoURL(input.VideoURL)
	if err != nil {
		return nil, err
	}
	input.VideoURL = cleanURL
	group, err := s.GRepo.FetchByPublicID(groupPublicID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	isMember, err := s.canManage(group.InternalID, authorID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you are not a member of this group")
	}
	input.PublicID = uuid.New()
	input.GroupInternalID = group.InternalID
	input.AuthorInternalID = authorID
	images := make([]models.DocumentationImage, 0, len(imagePaths))
	for i, path := range imagePaths {
		caption := ""
		if i < len(imageCaptions) {
			caption = imageCaptions[i]
		}
		images = append(images, models.DocumentationImage{
			PublicID: uuid.New(),
			FilePath: path,
			Caption:  caption,
		})
	}
	if err := s.GDRepo.Create(&input, images); err != nil {
		return nil, err
	}
	input.Images = images
	return &input, nil
}

func (s *GroupDocumentationServices) FetchByPublicID(publicID string) (*models.GroupDocumentation, error) {
	doc, err := s.GDRepo.FetchByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("documentation not found")
	}
	return doc, nil
}

func (s *GroupDocumentationServices) FetchByGroup(groupPublicID string, offset, limit int) ([]models.GroupDocumentation, int64, error) {
	group, err := s.GRepo.FetchByPublicID(groupPublicID)
	if err != nil {
		return nil, 0, err
	}
	if group == nil {
		return nil, 0, errors.New("group not found")
	}
	return s.GDRepo.FetchByGroup(group.InternalID, offset, limit)
}

func (s *GroupDocumentationServices) Update(authorID int64, publicID string, input models.GroupDocumentation) error {
	doc, err := s.GDRepo.FetchByPublicID(publicID)
	if err != nil {
		return err
	}
	if doc == nil {
		return errors.New("documentation not found")
	}
	isMember, err := s.canManage(doc.GroupInternalID, authorID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this group")
	}
	cleanURL, err := utils.NormalizeVideoURL(input.VideoURL)
	if err != nil {
		return err
	}
	doc.Title = input.Title
	doc.Narrative = input.Narrative
	doc.Location = input.Location
	doc.VideoURL = cleanURL
	doc.EventDate = input.EventDate
	return s.GDRepo.Update(doc)
}

func (s *GroupDocumentationServices) Delete(authorID int64, publicID string) error {
	doc, err := s.GDRepo.FetchByPublicID(publicID)
	if err != nil {
		return err
	}
	if doc == nil {
		return errors.New("documentation not found")
	}
	isMember, err := s.canManage(doc.GroupInternalID, authorID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this group")
	}
	return s.GDRepo.SoftDelete(publicID)
}

func (s *GroupDocumentationServices) AddImage(authorID int64, docPublicID string, filePath string, caption string) (*models.DocumentationImage, error) {
	doc, err := s.GDRepo.FetchByPublicID(docPublicID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("documentation not found")
	}
	isMember, err := s.canManage(doc.GroupInternalID, authorID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you are not a member of this group")
	}
	if len(doc.Images) >= utils.MaxImagesPerDocumentation {
		return nil, fmt.Errorf("maximum %d images per documentation", utils.MaxImagesPerDocumentation)
	}
	image := &models.DocumentationImage{
		PublicID:                uuid.New(),
		DocumentationInternalID: doc.InternalID,
		FilePath:                filePath,
		Caption:                 caption,
		DisplayOrder:            len(doc.Images),
	}
	if err := s.GDRepo.AddImage(image); err != nil {
		return nil, err
	}
	return image, nil
}

func (s *GroupDocumentationServices) DeleteImage(authorID int64, docPublicID string, imagePublicID string) error {
	doc, err := s.GDRepo.FetchByPublicID(docPublicID)
	if err != nil {
		return err
	}
	if doc == nil {
		return errors.New("documentation not found")
	}
	isMember, err := s.canManage(doc.GroupInternalID, authorID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this group")
	}
	return s.GDRepo.DeleteImage(imagePublicID)
}
