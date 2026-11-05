package repositories

import (
	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
)

type SubmissionCommentRepository interface {
	Create(comment *models.SubmissionComment) error
	FetchByPublicID(PublicID string) (*models.SubmissionComment, error)
	FetchBySubmission(submissionID int64) ([]models.SubmissionComment, error)
	FetchLatestRound(SubmissionID int64) (int, error)
	Delete(publicID uuid.UUID, authorInternalID int64) error
}

type SubmissionCommentRepositorys struct {
}

func NewSubmissionCommentRepository() SubmissionCommentRepository {
	return &SubmissionCommentRepositorys{}
}

func (r *SubmissionCommentRepositorys) Create(comment *models.SubmissionComment) error {
	return config.DB.Create(comment).Error
}

func (r *SubmissionCommentRepositorys) FetchByPublicID(PublicID string) (*models.SubmissionComment, error) {
	var comment models.SubmissionComment
	err := config.DB.Where("public_id  = ?", PublicID).First(&comment).Error
	return &comment, err
}

func (r *SubmissionCommentRepositorys) FetchBySubmission(submissionID int64) ([]models.SubmissionComment, error) {
	var comment []models.SubmissionComment
	err := config.DB.Preload("Author").Where("submission_internal_id", submissionID).Order("created_at ASC").Find(&comment).Error
	return comment, err
}

func (r *SubmissionCommentRepositorys) FetchLatestRound(SubmissionID int64) (int, error) {
	var RecentRound int
	err := config.DB.Model(&models.SubmissionComment{}).Where("submission_internal_id = ? AND decision IS NOT NULL", SubmissionID).Select("COALESCE(MAX(round), 0)").Scan(&RecentRound).Error
	return RecentRound, err
}
func (r *SubmissionCommentRepositorys) Delete(publicID uuid.UUID, authorInternalID int64) error {
	return config.DB.
		Where("public_id = ? AND author_internal_id = ?", publicID, authorInternalID).
		Delete(&models.SubmissionComment{}).Error
}
