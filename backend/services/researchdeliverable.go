package services

import (
	"errors"
	"strings"

	"time"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"gorm.io/gorm"
)

type ResearchDeliverableService interface {
	Create(userID int64, groupPublicID string, input models.ResearchDeliverable, filePath *string) (*models.ResearchDeliverable, error)
	FetchByGroup(groupPublicID string) ([]models.ResearchDeliverable, models.DeliverableProgress, error)
	Update(userID int64, deliverablePublicID string, input models.ResearchDeliverable) error
	MarkDone(userID int64, deliverablePublicID string, filePath *string) error
	MarkPending(userID int64, deliverablePublicID string) error
	Delete(userID int64, deliverablePublicID string) error
}

type DeliverableProgress struct {
	Total   int `json:"total"`
	Done    int `json:"done"`
	Pending int `json:"pending"`
	Percent int `json:"percent"`
}

type ResearchDeliverableServices struct {
	repo      repositories.ResearchDeliverableRepository
	groupRepo repositories.ResearchGroupRepository
	db        *gorm.DB
}

func NewResearchDeliverableService(repo repositories.ResearchDeliverableRepository, groupRepo repositories.ResearchGroupRepository, db *gorm.DB) ResearchDeliverableService {
	return &ResearchDeliverableServices{repo: repo, groupRepo: groupRepo, db: db}
}

func isValidProofType(t string) bool {
	return t == models.DeliverableProofCheckbox || t == models.DeliverableProofFile
}

func (s *ResearchDeliverableServices) Create(userID int64, groupPublicID string, input models.ResearchDeliverable, filePath *string) (*models.ResearchDeliverable, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title is required")
	}
	if input.ProofType == "" {
		input.ProofType = models.DeliverableProofCheckbox
	}
	if !isValidProofType(input.ProofType) {
		return nil, errors.New("invalid proof_type")
	}
	group, err := s.groupRepo.FetchByPublicID(groupPublicID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	isMember, err := s.groupRepo.IsMember(group.InternalID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you are not a member of this group")
	}
	order, err := s.repo.NextDisplayOrder(group.ResearchInternalID)
	if err != nil {
		return nil, err
	}
	d := &models.ResearchDeliverable{
		PublicID:            uuid.New(),
		ResearchInternalID:  group.ResearchInternalID,
		Title:               input.Title,
		Description:         input.Description,
		ProofType:           input.ProofType,
		Status:              models.DeliverableStatusPending,
		DisplayOrder:        order,
		CreatedByInternalID: userID,
	}
	if d.ProofType == models.DeliverableProofCheckbox {
		d.FilePath = nil
	} else if filePath != nil {
		d.FilePath = filePath
	}
	if err := s.repo.Create(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *ResearchDeliverableServices) FetchByGroup(groupPublicID string) ([]models.ResearchDeliverable, models.DeliverableProgress, error) {
	group, err := s.groupRepo.FetchByPublicID(groupPublicID)
	if err != nil {
		return nil, models.DeliverableProgress{}, err
	}
	if group == nil {
		return nil, models.DeliverableProgress{}, errors.New("group not found")
	}
	list, err := s.repo.FetchByResearch(group.ResearchInternalID)
	if err != nil {
		return nil, models.DeliverableProgress{}, err
	}
	total, done, err := s.repo.CountByStatus(group.ResearchInternalID)
	if err != nil {
		return nil, models.DeliverableProgress{}, err
	}
	percent := 0
	if total > 0 {
		percent = int((done * 100) / total)
	}
	progress := models.DeliverableProgress{
		Total:   int(total),
		Done:    int(done),
		Pending: int(total - done),
		Percent: percent,
	}
	return list, progress, nil
}

func (s *ResearchDeliverableServices) Update(userID int64, deliverablePublicID string, input models.ResearchDeliverable) error {
	d, err := s.repo.FetchByPublicID(deliverablePublicID)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("deliverable not found")
	}
	isMember, err := s.groupRepo.IsMember(d.ResearchInternalID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this group")
	}
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("title is required")
	}
	if input.ProofType != "" {
		if !isValidProofType(input.ProofType) {
			return errors.New("invalid proof_type")
		}
		d.ProofType = input.ProofType
	}
	d.Title = input.Title
	d.Description = input.Description
	d.DisplayOrder = input.DisplayOrder
	return s.repo.Update(d)
}

func (s *ResearchDeliverableServices) MarkDone(userID int64, deliverablePublicID string, filePath *string) error {
	d, err := s.repo.FetchByPublicID(deliverablePublicID)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("deliverable not found")
	}
	// isMember, err := s.groupRepo.IsMember(d.ResearchInternalID, userID)
	// if err != nil {
	// 	return err
	// }
	// if !isMember {
	// 	return errors.New("you are not a member of this group")
	// }
	if d.ProofType == models.DeliverableProofFile && (filePath == nil || *filePath == "") && d.FilePath == nil {
		return errors.New("a file is required to mark this deliverable as done")
	}
	now := time.Now()
	d.Status = models.DeliverableStatusDone
	d.CompletedByInternalID = &userID
	d.CompletedAt = &now
	if filePath != nil && *filePath != "" {
		d.FilePath = filePath
	}
	return s.repo.Update(d)
}

func (s *ResearchDeliverableServices) MarkPending(userID int64, deliverablePublicID string) error {
	d, err := s.repo.FetchByPublicID(deliverablePublicID)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("deliverable not found")
	}
	// isMember, err := s.groupRepo.IsMember(d.ResearchInternalID, userID)
	// if err != nil {
	// 	return err
	// }
	// if !isMember {
	// 	return errors.New("you are not a member of this group")
	// }
	d.Status = models.DeliverableStatusPending
	d.CompletedByInternalID = nil
	d.CompletedAt = nil
	return s.repo.Update(d)
}

func (s *ResearchDeliverableServices) Delete(userID int64, deliverablePublicID string) error {
	d, err := s.repo.FetchByPublicID(deliverablePublicID)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("deliverable not found")
	}
	isMember, err := s.groupRepo.IsMember(d.ResearchInternalID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this group")
	}
	return s.repo.SoftDelete(deliverablePublicID)
}
