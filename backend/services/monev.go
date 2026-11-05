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

type MonevService interface {
	CreateBatch(userID int64, input models.MonevBatch) (*models.MonevBatch, error)
	FetchBatchByPublicID(publicID string) (*models.MonevBatch, error)
	ListBatches(filter, sort string, offset, limit int) ([]models.MonevBatch, int64, error)
	UpdateBatch(userID int64, publicID string, input models.MonevBatch) error
	DeleteBatch(userID int64, publicID string) error

	AddOpdToBatch(userID int64, batchPublicID, opdPublicID string) (*models.MonevOpdEntry, error)
	RemoveOpdFromBatch(userID int64, entryPublicID string) error
	ListBatchEntries(userID int64, batchPublicID string) ([]models.MonevOpdEntry, error)

	FetchEntryByPublicID(userID int64, role string, entryPublicID string) (*models.MonevOpdEntry, error)
	ListMyEntries(userID int64, opdID int64) ([]models.MonevOpdEntry, error)
	FillRepresentative(userID int64, entryPublicID string, name, nip string) error
	SubmitEntry(userID int64, entryPublicID string) error
	VerifyEntry(userID int64, entryPublicID string) error

	AddKajian(userID int64, entryPublicID string, input models.KajianResearch) (*models.KajianResearch, error)
	UpdateKajian(userID int64, kajianPublicID string, input models.KajianResearch) error
	DeleteKajian(userID int64, kajianPublicID string) error

	AddPoint(userID int64, kajianPublicID string, input models.ResearchPoint) (*models.ResearchPoint, error)
	UpdatePoint(userID int64, role string, pointPublicID string, input models.ResearchPoint) error
	DeletePoint(userID int64, pointPublicID string) error
	GetUserOpd(userID int64) (*models.User, error)
}

type MonevServices struct {
	repo     repositories.MonevRepository
	userRepo repositories.UserRepository
	opdRepo  repositories.OpdRepository
	db       *gorm.DB
}

func NewMonevService(repo repositories.MonevRepository, userRepo repositories.UserRepository, opdRepo repositories.OpdRepository, db *gorm.DB) MonevService {
	return &MonevServices{repo: repo, userRepo: userRepo, opdRepo: opdRepo, db: db}
}

func (s *MonevServices) CreateBatch(userID int64, input models.MonevBatch) (*models.MonevBatch, error) {
	if input.Year == 0 {
		return nil, errors.New("year is required")
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title is required")
	}
	if input.Status == "" {
		input.Status = models.MonevStatusDraft
	}
	input.PublicID = uuid.New()
	input.CreatedByInternalID = userID
	if err := s.repo.CreateBatch(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *MonevServices) FetchBatchByPublicID(publicID string) (*models.MonevBatch, error) {
	b, err := s.repo.FetchBatchByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("batch not found")
	}
	return b, nil
}

func (s *MonevServices) ListBatches(filter, sort string, offset, limit int) ([]models.MonevBatch, int64, error) {
	return s.repo.ListBatches(filter, sort, offset, limit)
}

func (s *MonevServices) UpdateBatch(userID int64, publicID string, input models.MonevBatch) error {
	b, err := s.repo.FetchBatchByPublicID(publicID)
	if err != nil || b == nil {
		return errors.New("batch not found")
	}
	b.Year = input.Year
	b.Title = input.Title
	b.Description = input.Description
	if input.Status != "" {
		b.Status = input.Status
	}
	return s.repo.UpdateBatch(b)
}

func (s *MonevServices) DeleteBatch(userID int64, publicID string) error {
	b, err := s.repo.FetchBatchByPublicID(publicID)
	if err != nil || b == nil {
		return errors.New("batch not found")
	}
	return s.repo.SoftDeleteBatch(publicID)
}

func (s *MonevServices) AddOpdToBatch(userID int64, batchPublicID, opdPublicID string) (*models.MonevOpdEntry, error) {
	batch, err := s.repo.FetchBatchByPublicID(batchPublicID)
	if err != nil || batch == nil {
		return nil, errors.New("batch not found")
	}
	opd, err := s.opdRepo.FetchByPublicID(opdPublicID)
	if err != nil || opd == nil {
		return nil, errors.New("opd not found")
	}
	entry := &models.MonevOpdEntry{
		PublicID:        uuid.New(),
		BatchInternalID: batch.InternalID,
		OpdInternalID:   opd.InternalID,
		Status:          models.OpdEntryStatusPending,
	}
	if err := s.repo.CreateEntry(entry); err != nil {
		return nil, err
	}
	entry.Opd = opd
	return entry, nil
}

func (s *MonevServices) RemoveOpdFromBatch(userID int64, entryPublicID string) error {
	e, err := s.repo.FetchEntryByPublicID(entryPublicID)
	if err != nil || e == nil {
		return errors.New("entry not found")
	}
	return s.repo.SoftDeleteEntry(entryPublicID)
}

func (s *MonevServices) ListBatchEntries(userID int64, batchPublicID string) ([]models.MonevOpdEntry, error) {
	batch, err := s.repo.FetchBatchByPublicID(batchPublicID)
	if err != nil || batch == nil {
		return nil, errors.New("batch not found")
	}
	return s.repo.ListEntriesByBatch(batch.InternalID)
}

func (s *MonevServices) FetchEntryByPublicID(userID int64, role string, entryPublicID string) (*models.MonevOpdEntry, error) {
	e, err := s.repo.FetchEntryByPublicID(entryPublicID)
	if err != nil || e == nil {
		return nil, errors.New("entry not found")
	}
	if role == "opd" {
		user, err := s.userRepo.FetchByID(uint(userID))
		if err != nil || user == nil || user.OpdInternalID == nil || *user.OpdInternalID != e.OpdInternalID {
			return nil, errors.New("you do not have access to this entry")
		}
	}
	return e, nil
}

func (s *MonevServices) ListMyEntries(userID int64, opdID int64) ([]models.MonevOpdEntry, error) {
	return s.repo.ListEntriesByOpd(opdID)
}

func (s *MonevServices) FillRepresentative(userID int64, entryPublicID string, name, nip string) error {
	e, err := s.repo.FetchEntryByPublicID(entryPublicID)
	if err != nil || e == nil {
		return errors.New("entry not found")
	}
	user, err := s.userRepo.FetchByID(uint(userID))
	if err != nil || user == nil || user.OpdInternalID == nil || *user.OpdInternalID != e.OpdInternalID {
		return errors.New("you do not have access to this entry")
	}
	if e.Status == models.OpdEntryStatusVerified {
		return errors.New("entry already verified, cannot modify")
	}
	e.RepresentativeName = name
	e.RepresentativeNip = nip
	e.FilledByInternalID = &userID
	return s.repo.UpdateEntry(e)
}

func (s *MonevServices) SubmitEntry(userID int64, entryPublicID string) error {
	e, err := s.repo.FetchEntryByPublicID(entryPublicID)
	if err != nil || e == nil {
		return errors.New("entry not found")
	}
	user, err := s.userRepo.FetchByID(uint(userID))
	if err != nil || user == nil || user.OpdInternalID == nil || *user.OpdInternalID != e.OpdInternalID {
		return errors.New("you do not have access to this entry")
	}
	if e.RepresentativeName == "" || e.RepresentativeNip == "" {
		return errors.New("representative name and nip are required before submitting")
	}
	now := time.Now()
	e.Status = models.OpdEntryStatusSubmitted
	e.SubmittedAt = &now
	return s.repo.UpdateEntry(e)
}

func (s *MonevServices) VerifyEntry(userID int64, entryPublicID string) error {
	e, err := s.repo.FetchEntryByPublicID(entryPublicID)
	if err != nil || e == nil {
		return errors.New("entry not found")
	}
	now := time.Now()
	e.Status = models.OpdEntryStatusVerified
	e.VerifiedAt = &now
	return s.repo.UpdateEntry(e)
}

func (s *MonevServices) AddKajian(userID int64, entryPublicID string, input models.KajianResearch) (*models.KajianResearch, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title is required")
	}
	e, err := s.repo.FetchEntryByPublicID(entryPublicID)
	if err != nil || e == nil {
		return nil, errors.New("entry not found")
	}
	if e.Status == models.OpdEntryStatusVerified {
		return nil, errors.New("entry already verified, cannot modify")
	}
	k := &models.KajianResearch{
		PublicID:           uuid.New(),
		OpdEntryInternalID: e.InternalID,
		Title:              input.Title,
		Description:        input.Description,
		DisplayOrder:       input.DisplayOrder,
	}
	if err := s.repo.CreateKajian(k); err != nil {
		return nil, err
	}
	return k, nil
}

func (s *MonevServices) UpdateKajian(userID int64, kajianPublicID string, input models.KajianResearch) error {
	k, err := s.repo.FetchKajianByPublicID(kajianPublicID)
	if err != nil || k == nil {
		return errors.New("kajian not found")
	}
	if k.OpdEntry != nil && k.OpdEntry.Status == models.OpdEntryStatusVerified {
		return errors.New("entry already verified, cannot modify")
	}
	k.Title = input.Title
	k.Description = input.Description
	k.DisplayOrder = input.DisplayOrder
	return s.repo.UpdateKajian(k)
}

func (s *MonevServices) DeleteKajian(userID int64, kajianPublicID string) error {
	k, err := s.repo.FetchKajianByPublicID(kajianPublicID)
	if err != nil || k == nil {
		return errors.New("kajian not found")
	}
	return s.repo.SoftDeleteKajian(kajianPublicID)
}

func (s *MonevServices) AddPoint(userID int64, kajianPublicID string, input models.ResearchPoint) (*models.ResearchPoint, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title is required")
	}
	k, err := s.repo.FetchKajianByPublicID(kajianPublicID)
	if err != nil || k == nil {
		return nil, errors.New("kajian not found")
	}
	order, err := s.repo.NextPointOrder(k.InternalID)
	if err != nil {
		return nil, err
	}
	p := &models.ResearchPoint{
		PublicID:         uuid.New(),
		KajianInternalID: k.InternalID,
		Title:            input.Title,
		DisplayOrder:     order,
	}
	if err := s.repo.CreatePoint(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *MonevServices) UpdatePoint(userID int64, role string, pointPublicID string, input models.ResearchPoint) error {
	p, err := s.repo.FetchPointByPublicID(pointPublicID)
	if err != nil || p == nil {
		return errors.New("point not found")
	}
	// Admin: bisa update semua field termasuk title
	// OPD: hanya isi field monitoring, tidak bisa ganti title
	if role == "admin" {
		if strings.TrimSpace(input.Title) != "" {
			p.Title = input.Title
		}
	}
	p.MonitoringFilePath = input.MonitoringFilePath
	p.HasMonitoringResult = input.HasMonitoringResult
	p.MonitoringResultDescription = input.MonitoringResultDescription
	p.KendalaPelaksanaanRekomendasi = input.KendalaPelaksanaanRekomendasi
	p.Benefits = input.Benefits
	return s.repo.UpdatePoint(p)
}

func (s *MonevServices) DeletePoint(userID int64, pointPublicID string) error {
	p, err := s.repo.FetchPointByPublicID(pointPublicID)
	if err != nil || p == nil {
		return errors.New("point not found")
	}
	return s.repo.SoftDeletePoint(pointPublicID)
}
func (s *MonevServices) GetUserOpd(userID int64) (*models.User, error) {
	return s.userRepo.FetchByID(uint(userID))
}
