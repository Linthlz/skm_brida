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

type ResearchService interface {
	Submit(uploaderPublicID *string, research *models.Research) error
	Revision(researcherID int64, publicID string, req models.UpdateSubmission) error
	ResearchSlug(slug string) (*models.Research, error)
	FetchByPublicID(publicID string) (*models.Research, error)
	FetchByID(id uint) (*models.Research, error)
	Researches(filter, sort string, limit, offset int) ([]models.Research, int64, error)
	UserResearches(userInternalID uint, filter, sort string, limit, offset int) ([]models.Research, int64, error)
}

type ResearchServices struct {
	userRepo     repositories.UserRepository
	researchRepo repositories.ResearchRepository
	NotifRepo    repositories.NotificationRepository
	notifService NotificationService
}

func uniqueResearchSlug(tx *gorm.DB, title string, excludeID int64) (string, error) {
	base := utils.SlugifyResearch(title)
	candidate := base
	for i := 2; i <= 50; i++ {
		var count int64
		q := tx.Session(&gorm.Session{NewDB: true}).Model(&models.Research{}).Where("slug = ?", candidate)
		if excludeID > 0 {
			q = q.Where("research_id <> ?", excludeID)
		}
		if err := q.Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return "", utils.ErrConflict("too many submissions with a similar title")
}

func NewResearchService(userRepo repositories.UserRepository, researchRepo repositories.ResearchRepository, notifRepo repositories.NotificationRepository, notifService NotificationService) ResearchService {
	return &ResearchServices{userRepo, researchRepo, notifRepo, notifService}
}

func (s *ResearchServices) Submit(uploaderPublicID *string, research *models.Research) error {
	user, err := s.userRepo.FetchByPublicID(*uploaderPublicID)
	if err != nil {
		return errors.New("user not found")
	}
	if research.Title == nil {
		return errors.New("Title Required!")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		research.PublicID = uuid.New()
		research.UploaderInternalID = &user.InternalID
		research.Status = models.SubmissionPending
		slug, err := uniqueResearchSlug(tx, *research.Title, 0)
		if err != nil {
			return err
		}
		research.Slug = slug
		if err := tx.Create(research).Error; err != nil {
			return err
		}
		year := time.Now().Year()
		code := fmt.Sprintf("BRD%dRSH-%04d", year, *research.InternalID)
		research.ResearchCode = &code
		if err := tx.Model(models.Research{}).Where("research_id = ?", research.InternalID).Update("research_code", code).Error; err != nil {
			return err
		}
		if err := s.notifService.ReviewSubmissionNotification(research, user); err != nil {
			log.Printf("[+] failed to notify reviewers: %v", err)
		}
		return nil
	})
}

//	func (s *ResearchServices) Revision(PublicID string, req models.UpdateSubmission) error {
//		updates := map[string]interface{}{}
//		if req.Title != nil {
//			updates["title"] = req.Title
//		}
//		if req.Purpose != nil {
//			updates["purpose"] = req.Purpose
//		}
//		if req.Category != nil {
//			updates["category"] = req.Category
//		}
//		if req.Funding_Scheme != nil {
//			updates["funding_scheme"] = req.Funding_Scheme
//		}
//		if req.Fund_Amount != nil {
//			updates["fund_amount"] = req.Fund_Amount
//		}
//		if req.Location != nil {
//			updates["location"] = req.Location
//		}
//		if req.Address != nil {
//			updates["address"] = req.Address
//		}
//		if req.RPJMD != nil {
//			updates["rpjmd"] = req.RPJMD
//		}
//		if req.Promised_Output != nil {
//			updates["promised_output"] = req.Promised_Output
//		}
//		if req.Research_Significance != nil {
//			updates["research_significance"] = req.Research_Significance
//		}
//		if req.Research_Time_Range != nil {
//			updates["research_time_range"] = req.Research_Time_Range
//		}
//		if req.Partnership_Target_Selection != nil {
//			updates["partnership_target_selection"] = req.Partnership_Target_Selection
//		}
//		if req.Partner_Name != nil {
//			updates["partner_name"] = req.Partner_Name
//		}
//		if req.Estimated_Number_of_Direct_Beneficiaries != nil {
//			updates["estimated_number_of_direct_beneficiaries"] = req.Estimated_Number_of_Direct_Beneficiaries
//		}
//		if req.File != nil {
//			updates["research_proposal_file"] = req.File
//		}
//		return s.researchRepo.SubmissionUpdate(PublicID, updates)
//	}
func (s *ResearchServices) Revision(researcherID int64, publicID string, req models.UpdateSubmission) error {
	submission, err := s.researchRepo.FetchByPublicID(publicID)
	if err != nil {
		return err
	}
	if submission == nil {
		return errors.New("submission not found")
	}
	if submission.UploaderInternalID == nil || *submission.UploaderInternalID != researcherID {
		return utils.ErrForbidden("you are not the owner of this submission")
	}
	if submission.Status != models.SubmissionRevision {
		return utils.ErrValidation("only submissions with status 'revision' can be revised")
	}

	updates := map[string]interface{}{
		"status": models.SubmissionPending,
	}
	if req.Title != nil && strings.TrimSpace(*req.Title) != "" && (submission.Title == nil || *req.Title != *submission.Title) {
		slug, err := uniqueResearchSlug(config.DB, *req.Title, *submission.InternalID)
		if err != nil {
			return err
		}
		updates["title"] = strings.TrimSpace(*req.Title)
		updates["slug"] = slug
	}
	if req.Purpose != nil {
		updates["purpose"] = *req.Purpose
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Funding_Scheme != nil {
		updates["funding_scheme"] = *req.Funding_Scheme
	}
	if req.Fund_Amount != nil {
		updates["fund_amount"] = *req.Fund_Amount
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.RPJMD != nil {
		updates["rpjmd"] = *req.RPJMD
	}
	if req.Promised_Output != nil {
		updates["promised_output"] = *req.Promised_Output
	}
	if req.Research_Significance != nil {
		updates["research_significance"] = *req.Research_Significance
	}
	if req.Research_Time_Range != nil {
		updates["research_time_range"] = *req.Research_Time_Range
	}
	if req.Partnership_Target_Selection != nil {
		updates["partnership_target_selection"] = *req.Partnership_Target_Selection
	}
	if req.Partner_Name != nil {
		updates["partner_name"] = *req.Partner_Name
	}
	if req.Estimated_Number_of_Direct_Beneficiaries != nil {
		updates["estimated_number_of_direct_beneficiaries"] = *req.Estimated_Number_of_Direct_Beneficiaries
	}
	// research_proposal_file is never taken from the request body: clients
	// must not be able to point the record at an arbitrary server path.

	if err := s.researchRepo.SubmissionUpdate(publicID, updates); err != nil {
		return err
	}

	// if submission.ReviewerInternalID != nil {
	// 	go s.notifService.NotifyUser(
	// 		*submission.ReviewerInternalID,
	// 		"Researcher Submitted a Revision",
	// 		fmt.Sprintf(`A revised version of "%s" is ready for review.`, safeTitle(submission)),
	// 		"submission",
	// 		submission.InternalID,
	// 		submission.PublicID,
	// 	)
	// }

	return nil
}
func (s *ResearchServices) FetchByPublicID(publicID string) (*models.Research, error) {
	return s.researchRepo.FetchByPublicID(publicID)
}

func (s *ResearchServices) FetchByID(id uint) (*models.Research, error) {
	return s.researchRepo.FetchByID(id)
}
func (s *ResearchServices) ResearchSlug(slug string) (*models.Research, error) {
	return s.researchRepo.ResearchSlug(slug)
}
func (s *ResearchServices) Researches(filter, sort string, limit, offset int) ([]models.Research, int64, error) {
	return s.researchRepo.Researches(filter, sort, limit, offset)
}

func (s *ResearchServices) UserResearches(userInternalID uint, filter, sort string, limit, offset int) ([]models.Research, int64, error) {
	return s.researchRepo.UserResearches(userInternalID, filter, sort, limit, offset)
}
