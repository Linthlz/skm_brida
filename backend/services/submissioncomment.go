package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"gorm.io/gorm"
)

type SubmissionCommentService interface {
	AssignReviewer(ReviewerID int64, SubmissionPublicID string) error
	AddReview(ReviewerID int64, SubmissionPublicID string, decision models.SubmissionDecision, message string) (*models.SubmissionComment, error)
	ApproveSubmission(ReviewerID int64, SubmissionPublicID string, message string) error
}
type SubmissionCommentServices struct {
	repo         repositories.SubmissionCommentRepository
	researchRepo repositories.ResearchRepository
	userRepo     repositories.UserRepository
	notifService NotificationService
	db           *gorm.DB
	groupService ResearchGroupService
}

func NewSubmissionCommentService(repo repositories.SubmissionCommentRepository, researchRepo repositories.ResearchRepository, userRepo repositories.UserRepository, notifService NotificationService, db *gorm.DB, groupService ResearchGroupService) SubmissionCommentService {
	return &SubmissionCommentServices{repo: repo, researchRepo: researchRepo, userRepo: userRepo, notifService: notifService, db: db, groupService: groupService}
}

func (s *SubmissionCommentServices) AssignReviewer(ReviewerID int64, SubmissionPublicID string) error {
	submission, err := s.researchRepo.FetchByPublicID(SubmissionPublicID)
	if err != nil {
		return fmt.Errorf("Submission Not Found")
	}
	if submission.UploaderInternalID != nil && *submission.UploaderInternalID == ReviewerID {
		return errors.New("You are forbidden to review your own submission")
	}
	// A revised submission returns to "pending" but keeps its reviewer, who may claim it again.
	claim := config.DB.Model(&models.Research{}).Where("public_id = ? AND status = ? AND (reviewer_internal_id IS NULL OR reviewer_internal_id = ?)", SubmissionPublicID, models.SubmissionPending, ReviewerID).Updates(map[string]interface{}{
		"reviewer_internal_id": ReviewerID,
		"status":               models.SubmissionUnderReview,
		"review_started_at":    time.Now(),
	})
	if claim.Error != nil {
		return claim.Error
	}
	if claim.RowsAffected == 0 {
		return errors.New("Submission already claimed!")
	}
	return nil
}

func (s *SubmissionCommentServices) ApproveSubmission(ReviewerID int64, SubmissionPublicID string, message string) error {
	submission, err := s.researchRepo.FetchByPublicID(SubmissionPublicID)
	if err != nil {
		return err
	}
	if submission == nil {
		return errors.New("submission not found")
	}
	if submission.ReviewerInternalID == nil || *submission.ReviewerInternalID != ReviewerID {
		return errors.New("you are not assigned to this submission")
	}
	if submission.Status != models.SubmissionUnderReview {
		return fmt.Errorf("cannot approve — status is %s", submission.Status)
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var lastRound int
		if err := tx.Model(&models.SubmissionComment{}).
			Where("submission_internal_id = ? AND decision IS NOT NULL", submission.InternalID).
			Select("COALESCE(MAX(round), 0)").
			Scan(&lastRound).Error; err != nil {
			return err
		}
		nextRound := lastRound + 1

		decision := models.DecisionApproved
		comment := &models.SubmissionComment{
			PublicID:             uuid.New(),
			SubmissionInternalID: *submission.InternalID,
			AuthorInternalID:     ReviewerID,
			Message:              message,
			Decision:             &decision,
			Round:                &nextRound,
		}
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		group, err := s.groupService.CreateGroup(tx, submission)
		if err != nil {
			return err
		}
		now := time.Now()
		return tx.Model(&models.Research{}).
			Where("research_id = ?", submission.InternalID).
			Select("status", "review_finished_at", "group_internal_id").
			Updates(map[string]interface{}{
				"status":             models.SubmissionOnGoing,
				"review_finished_at": now,
				"group_internal_id":  group.InternalID,
			}).Error
	})
}

func (s *SubmissionCommentServices) AddReview(ReviewerID int64, SubmissionPublicID string, decision models.SubmissionDecision, message string) (*models.SubmissionComment, error) {
	if message == "" {
		return nil, errors.New("comment cannot be empty")
	}
	switch decision {
	case models.DecisionRejected, models.DecisionRevision:
	case models.DecisionApproved:
		return nil, errors.New("use the approve endpoint to approve a submission")
	default:
		return nil, fmt.Errorf("invalid decision: %q", decision)
	}

	submission, err := s.researchRepo.FetchByPublicID(SubmissionPublicID)
	if err != nil {
		return nil, err
	}
	if submission == nil {
		return nil, errors.New("submission not found")
	}
	if submission.Status != models.SubmissionUnderReview {
		return nil, fmt.Errorf("cannot submit decision, status is %s", submission.Status)
	}
	reviewer, err := s.userRepo.FetchByID(uint(ReviewerID))
	if err != nil || reviewer == nil {
		return nil, errors.New("reviewer not found")
	}

	var result *models.SubmissionComment
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var lastRound int
		if err := tx.Model(&models.SubmissionComment{}).Where("submission_internal_id = ? AND decision IS NOT NULL", submission.InternalID).Select("COALESCE(MAX(round), 0)").Scan(&lastRound).Error; err != nil {
			return err
		}
		nextRound := lastRound + 1
		comment := &models.SubmissionComment{
			PublicID:             uuid.New(),
			SubmissionInternalID: *submission.InternalID,
			AuthorInternalID:     reviewer.InternalID,
			Message:              message,
			Decision:             &decision,
			Round:                &nextRound,
		}

		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		if submission.ReviewerInternalID == nil || *submission.ReviewerInternalID != ReviewerID {
			return errors.New("you are not assigned to this submission")
		}
		status := models.SubmissionRevision
		if decision == models.DecisionRejected {
			status = models.SubmissionRejected
		}
		if err := tx.Model(&models.Research{}).Where("research_id = ?", submission.InternalID).Updates(map[string]interface{}{
			"status":             status,
			"review_finished_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		result = comment

		return nil
	})
	if err != nil {
		return nil, err
	}
	go s.notifService.NotifyDecision(submission, decision, message)
	return result, nil

}
