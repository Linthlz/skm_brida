package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"gorm.io/gorm"
)

type ResearchGroupService interface {
	CreateGroup(tx *gorm.DB, submission *models.Research) (*models.ResearchGroup, error)
	BackfillMissingGroups() (int, error)
	FetchByPublicID(publicID string) (*models.ResearchGroup, error)
	FetchMembers(groupPublicID string) ([]models.ResearchGroupMember, error)
	JoinByCode(userID int64, code string) (*models.ResearchGroup, error)
	FetchBySlug(slug string) (*models.ResearchGroup, error)
	FetchAll(filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error)
	ListMyGroups(userID int64, filter, sort string, offset, limit int) ([]models.PublicGroupOutput, int64, error)
	GetMyGroup(userID int64, publicID string) (*models.MyGroupDetailOutput, error)
	//Public
	ListPublic(filter, sort string, offset, limit int) ([]models.PublicGroupOutput, int64, error)
	GetPublicBySlug(slug string) (*models.PublicGroupOutput, error)
	GetPublicByPublicID(publicID string) (*models.PublicGroupOutput, error)
}

type ResearchGroupServices struct {
	repo         repositories.ResearchGroupRepository
	userRepo     repositories.UserRepository
	notifService NotificationService
	db           *gorm.DB
}

func NewResearchGroupService(repo repositories.ResearchGroupRepository, notifService NotificationService, db *gorm.DB, userRepo repositories.UserRepository) ResearchGroupService {
	return &ResearchGroupServices{repo: repo, notifService: notifService, db: db, userRepo: userRepo}
}
// CreateGroup membuat kelompok untuk riset yang disetujui (ketua = pengusul).
// Bila kelompok untuk riset ini sudah ada, kelompok tersebut dipakai ulang.
func (s *ResearchGroupServices) CreateGroup(tx *gorm.DB, submission *models.Research) (*models.ResearchGroup, error) {
	if submission.UploaderInternalID == nil {
		return nil, errors.New("submission has no uploader")
	}
	if submission.InternalID == nil {
		return nil, errors.New("submission has no internal id")
	}
	var existing models.ResearchGroup
	err := tx.Where("research_internal_id = ?", *submission.InternalID).First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	title := "Untitled Research"
	if submission.Title != nil && *submission.Title != "" {
		title = *submission.Title
	}
	slug, err := uniqueGroupSlug(tx, submission.Slug, title)
	if err != nil {
		return nil, err
	}
	group := &models.ResearchGroup{
		PublicID:           uuid.New(),
		ResearchInternalID: *submission.InternalID,
		Name:               title,
		Slug:               &slug,
		JoinCode:           utils.GenerateJoinCode(title),
		LeaderInternalID:   *submission.UploaderInternalID,
		Description:        "Auto-created on approval",
	}
	if err := tx.Omit("Leader", "Members").Create(group).Error; err != nil {
		return nil, err
	}
	leader := &models.ResearchGroupMember{
		PublicID:        uuid.New(),
		GroupInternalID: group.InternalID,
		UserInternalID:  *submission.UploaderInternalID,
		Role:            "leader",
		JoinedAt:        time.Now(),
	}
	if err := tx.Omit("User").Create(leader).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// uniqueGroupSlug memakai slug riset (sudah unik) dan menambah akhiran bila bentrok.
func uniqueGroupSlug(tx *gorm.DB, researchSlug, title string) (string, error) {
	base := strings.TrimSpace(researchSlug)
	if base == "" {
		base = utils.SlugifyResearch(title)
	}
	slug := base
	for i := 2; i < 100; i++ {
		var count int64
		if err := tx.Unscoped().Model(&models.ResearchGroup{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
	return base + "-" + uuid.NewString()[:8], nil
}

// BackfillMissingGroups membuat kelompok untuk riset berstatus on-going/approved yang
// belum memiliki kelompok (mis. disetujui sebelum fitur kelompok tersedia).
func (s *ResearchGroupServices) BackfillMissingGroups() (int, error) {
	db := s.db
	if db == nil {
		db = config.DB
	}
	var list []models.Research
	if err := db.Where("status IN ? AND group_internal_id IS NULL AND deleted_at IS NULL",
		[]models.SubmissionStatus{models.SubmissionOnGoing, models.SubmissionApproved}).Find(&list).Error; err != nil {
		return 0, err
	}
	created := 0
	for i := range list {
		r := &list[i]
		err := db.Transaction(func(tx *gorm.DB) error {
			group, err := s.CreateGroup(tx, r)
			if err != nil {
				return err
			}
			return tx.Model(&models.Research{}).Where("research_id = ?", *r.InternalID).
				Update("group_internal_id", group.InternalID).Error
		})
		if err != nil {
			code := ""
			if r.ResearchCode != nil {
				code = *r.ResearchCode
			}
			log.Printf("[groups] backfill gagal untuk riset %s: %v", code, err)
			continue
		}
		created++
	}
	return created, nil
}

func (s *ResearchGroupServices) FetchByPublicID(publicID string) (*models.ResearchGroup, error) {
	group, err := s.repo.FetchByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	return group, nil
}

func (s *ResearchGroupServices) FetchMembers(groupPublicID string) ([]models.ResearchGroupMember, error) {
	group, err := s.repo.FetchByPublicID(groupPublicID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	return s.repo.FetchMembers(group.InternalID)
}
func (s *ResearchGroupServices) JoinByCode(userID int64, code string) (*models.ResearchGroup, error) {
	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	if cleanCode == "" {
		return nil, errors.New("join code is required")
	}
	group, err := s.repo.FetchByJoinCode(cleanCode)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("invalid join code")
	}
	isMember, err := s.repo.IsMember(group.InternalID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, errors.New("you are already member of this group")
	}
	member := &models.ResearchGroupMember{
		PublicID:        uuid.New(),
		GroupInternalID: group.InternalID,
		UserInternalID:  userID,
		Role:            "member",
	}
	if err := s.repo.AddMember(member); err != nil {
		return nil, err
	}
	newMemberUser, _ := s.userRepo.FetchByID(uint(userID))
	go s.notifService.NotifyMemberJoined(group, newMemberUser)
	return group, nil
}
func (s *ResearchGroupServices) FetchBySlug(slug string) (*models.ResearchGroup, error) {
	group, err := s.repo.FetchBySlug(slug)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	return group, nil
}

func (s *ResearchGroupServices) FetchAll(filter, sort string, offset, limit int) ([]models.ResearchGroup, int64, error) {
	return s.repo.FetchAll(filter, sort, offset, limit)
}
func toPublicGroup(group *models.ResearchGroup) models.PublicGroupOutput {
	out := models.PublicGroupOutput{
		PublicID:    group.PublicID,
		Name:        group.Name,
		Slug:        group.Slug,
		Description: group.Description,
		MemberCount: len(group.Members),
		CreatedAt:   group.CreatedAt,
	}

	if group.Leader != nil {
		out.Leader.PublicID = group.Leader.PublicID
		out.Leader.Position = group.Leader.Position
		if group.Leader.Name != nil {
			out.Leader.Name = *group.Leader.Name
		}
	}

	if group.Research != nil {
		out.Research = &models.PublicResearchOutput{
			PublicID:              group.Research.PublicID,
			Slug:                  &group.Research.Slug,
			Category:              group.Research.Category,
			Location:              group.Research.Location,
			Address:               group.Research.Address,
			Research_Time_Range:   group.Research.Research_Time_Range,
			Status:                string(group.Research.Status),
			Purpose:               group.Research.Purpose,
			Research_Significance: group.Research.Research_Significance,
			Promised_Output:       group.Research.Promised_Output,
		}
		if group.Research.Title != nil {
			out.Research.Title = *group.Research.Title
		}
	}

	// Fallback: if group has no description, use research purpose
	if out.Description == "" && out.Research != nil {
		out.Description = out.Research.Purpose
	}

	return out
}

func (s *ResearchGroupServices) ListMyGroups(userID int64, filter, sort string, offset, limit int) ([]models.PublicGroupOutput, int64, error) {
	groups, total, err := s.repo.FetchMyGroups(userID, filter, sort, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	response := make([]models.PublicGroupOutput, 0, len(groups))
	for i := range groups {
		response = append(response, toPublicGroup(&groups[i]))
	}
	return response, total, nil
}

func (s *ResearchGroupServices) GetMyGroup(userID int64, publicID string) (*models.MyGroupDetailOutput, error) {
	group, err := s.repo.FetchMyGroupByPublicID(userID, publicID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}

	out := &models.MyGroupDetailOutput{
		PublicID:    group.PublicID,
		Name:        group.Name,
		Slug:        group.Slug,
		Description: group.Description,
		JoinCode:    group.JoinCode,
		CreatedAt:   group.CreatedAt,
	}

	// My role
	if group.LeaderInternalID == userID {
		out.MyRole = "leader"
	} else {
		out.MyRole = "member"
	}

	// Leader
	if group.Leader != nil {
		out.Leader.PublicID = group.Leader.PublicID
		out.Leader.Position = group.Leader.Position
		if group.Leader.Name != nil {
			out.Leader.Name = *group.Leader.Name
		}
	}

	// Members — full list with names
	for _, m := range group.Members {
		var member struct {
			PublicID uuid.UUID `json:"public_id"`
			Name     string    `json:"name"`
			Position string    `json:"position"`
			Role     string    `json:"role"`
			JoinedAt time.Time `json:"joined_at"`
		}
		member.Role = m.Role
		member.JoinedAt = m.JoinedAt
		if m.User != nil {
			member.PublicID = m.User.PublicID
			member.Position = m.User.Position
			if m.User.Name != nil {
				member.Name = *m.User.Name
			}
		}
		out.Members = append(out.Members, member)
	}

	// Research
	if group.Research != nil {
		out.Research = &models.PublicResearchOutput{
			PublicID:              group.Research.PublicID,
			Slug:                  &group.Research.Slug,
			Category:              group.Research.Category,
			Location:              group.Research.Location,
			Address:               group.Research.Address,
			Research_Time_Range:   group.Research.Research_Time_Range,
			Status:                string(group.Research.Status),
			Purpose:               group.Research.Purpose,
			Research_Significance: group.Research.Research_Significance,
			Promised_Output:       group.Research.Promised_Output,
		}
		if group.Research.Title != nil {
			out.Research.Title = *group.Research.Title
		}
	}

	if out.Description == "" && out.Research != nil {
		out.Description = out.Research.Purpose
	}

	return out, nil
}

func (s *ResearchGroupServices) ListPublic(filter, sort string, offset, limit int) ([]models.PublicGroupOutput, int64, error) {
	groups, total, err := s.repo.ListPublic(filter, sort, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	response := make([]models.PublicGroupOutput, 0, len(groups))
	for i := range groups {
		response = append(response, toPublicGroup(&groups[i]))
	}
	return response, total, nil
}

func (s *ResearchGroupServices) GetPublicBySlug(slug string) (*models.PublicGroupOutput, error) {
	group, err := s.repo.FetchPublicBySlug(slug)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	out := toPublicGroup(group)
	return &out, nil
}

func (s *ResearchGroupServices) GetPublicByPublicID(publicID string) (*models.PublicGroupOutput, error) {
	group, err := s.repo.FetchPublicByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	out := toPublicGroup(group)
	return &out, nil
}
