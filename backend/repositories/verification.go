package repositories

import (
	"errors"
	"time"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/domain"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type VerificationRepository struct {
	verifications map[string]domain.Verification
}

func NewVerificationRepository() *VerificationRepository {
	return &VerificationRepository{}
}

func (r *VerificationRepository) Store(email string, verification domain.Verification) error {
	return config.DB.Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"verification_token":      verification.Code,
		"verification_expires_at": verification.Exp,
	}).Error
}

func (r *VerificationRepository) Fetch(email string) (domain.Verification, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Verification{}, err
	}
	if err != nil {
		return domain.Verification{}, err
	}
	if user.VerificationToken == nil || user.VerificationExpiry == nil {
		return domain.Verification{}, errors.New("token not found")
	}
	return domain.Verification{
		Code: *user.VerificationToken,
		Exp:  *user.VerificationExpiry,
	}, nil
}

func (r *VerificationRepository) Delete(email string) error {
	return config.DB.Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"verified":                true,
		"verified_at":             time.Now(),
		"verification_token":      nil,
		"verification_expires_at": nil,
	}).Error
}
