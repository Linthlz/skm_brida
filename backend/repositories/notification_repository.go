package repositories

import (
	"errors"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	NotificationPiles(RecipientID int64, filter, sort string, limit, offset int) ([]models.Notification, int64, error)
	FetchByPublicID(publicID string) (*models.Notification, error)
}

type NotificationRepositorys struct{}

func NewNotificationRepository() NotificationRepository {
	return &NotificationRepositorys{}
}

func (r *NotificationRepositorys) Create(notification *models.Notification) error {
	return config.DB.Create(notification).Error
}

func (r *NotificationRepositorys) NotificationPiles(RecipientID int64, filter, sort string, limit, offset int) ([]models.Notification, int64, error) {
	var pile []models.Notification
	var total int64
	db := config.DB.Model(&models.Notification{}).Where("recipient_internal_id = ?", RecipientID)
	if filter != "" {
		pattern := "%" + filter + "%"
		db = db.Where("subject ILIKE ? OR message ILIKE ?", pattern, pattern)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	switch sort {
	case "oldest":
		db = db.Order("created_at ASC")
	case "subject":
		db = db.Order("subject ASC")
	default:
		db = db.Order("created_at DESC")
	}
	err := db.Limit(limit).Offset(offset).Find(&pile).Error
	return pile, total, err
}

func (r *NotificationRepositorys) FetchByPublicID(publicID string) (*models.Notification, error) {
	var notification models.Notification
	err := config.DB.Where("public_id = ?", publicID).First(&notification).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &notification, nil
}
