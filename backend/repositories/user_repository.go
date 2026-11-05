package repositories

import (
	"time"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
)

type UserRepository interface {
	// Create
	Create(user *models.User) error
	//Update Delete
	// UserUpdate(user *models.User) error
	UserUpdate(publicID string, updates map[string]interface{}) error
	UserDelete(id uint) error
	//Read
	FetchByEmail(email string) (*models.User, error)
	FetchByPublicID(publicID string) (*models.User, error)
	FetchByID(id uint) (*models.User, error)
	FetchUsersPaginated(filter, role, sort string, limit, offset int) ([]models.User, int64, error)
	CountAdmins() (int64, error)
	UserDependencyCount(userID int64) (int64, error)
	//Email
	StoreToken(email string, token string, expired time.Time) error
	FetchUserVerificationToken(VerificationToken string) (*models.User, error)
	VerifiedEmail(publicID string) error
}

type UserRepositorys struct{}

func NewUserRepository() UserRepository {
	return &UserRepositorys{}
}

func (r *UserRepositorys) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *UserRepositorys) FetchByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user)
	return &user, err.Error
}

func (r *UserRepositorys) FetchByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id= ?", publicID).First(&user).Error
	return &user, err
}
func (r *UserRepositorys) FetchByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

//	func (r *UserRepositorys) UserUpdate(user *models.User) error {
//		return config.DB.Model(&models.User{}).Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
//			// "name":                user.Name,
//			// "phone_number":        user.Phone_Number,
//			// "institution":         user.Institution,
//			// "institution_address": user.Institution_Address,
//			// "position":            user.Position,
//		}).Error
//	}
func (r *UserRepositorys) UserUpdate(publicID string, updates map[string]interface{}) error {
	return config.DB.Model(&models.User{}).Where("public_id=?", publicID).Updates(updates).Error
}
func (r *UserRepositorys) UserDelete(id uint) error {
	return config.DB.Delete(models.User{}, id).Error
}

func (r *UserRepositorys) StoreToken(email string, token string, expired time.Time) error {
	return config.DB.Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"verification_token":      token,
		"verification_expires_at": expired,
	}).Error
}
func (r *UserRepositorys) FetchUserVerificationToken(VerificationToken string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("verification_token = ?", VerificationToken).First(&user).Error
	return &user, err
}

func (r *UserRepositorys) VerifiedEmail(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Updates(map[string]interface{}{
		"verified":                true,
		"verified_at":             time.Now(),
		"verification_token":      nil,
		"verification_expires_at": nil,
	}).Error
}

func (r *UserRepositorys) FetchUsersPaginated(filter, role, sort string, limit, offset int) ([]models.User, int64, error) {
	var user []models.User
	var total int64

	db := config.DB.Model(&models.User{})
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("(name ILIKE ? OR email ILIKE ?)", filterPattern, filterPattern)
	}
	if role != "" {
		db = db.Where("role = ?", role)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// if sort != "" {
	// 	if sort == "-id" {
	// 		sort = "-internal_id"
	// 	} else if sort == "id" {
	// 		sort = "internal_id "
	// 	}
	// 	if strings.HasPrefix(sort, "") {
	// 		sort = strings.TrimPrefix(sort, "-") + "DESC"
	// 	} else {
	// 		sort += " ASC"
	// 	}
	// 	db = db.Order(sort)
	// }
	switch sort {
	case "oldest":
		db = db.Order("created_at ASC")
	case "subject":
		db = db.Order("subject ASC")
	default:
		db = db.Order("created_at DESC")
	}
	err := db.Limit(limit).Offset(offset).Find(&user).Error
	return user, total, err
}

func (r *UserRepositorys) CountAdmins() (int64, error) {
	var n int64
	err := config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&n).Error
	return n, err
}

// UserDependencyCount counts rows that reference the user and would block deletion
// (their foreign keys are NOT NULL, so SET NULL cannot apply).
func (r *UserRepositorys) UserDependencyCount(userID int64) (int64, error) {
	var n int64
	err := config.DB.Raw(`
		SELECT
		  (SELECT COUNT(*) FROM researches WHERE uploader_internal_id = @id) +
		  (SELECT COUNT(*) FROM news WHERE author_internal_id = @id) +
		  (SELECT COUNT(*) FROM group_documentations WHERE author_internal_id = @id) +
		  (SELECT COUNT(*) FROM submission_comments WHERE author_internal_id = @id) +
		  (SELECT COUNT(*) FROM research_groups WHERE leader_internal_id = @id)`,
		map[string]interface{}{"id": userID}).Scan(&n).Error
	return n, err
}
