package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/infrastructure/code"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	FetchUserByPublicID(publicID string) (*models.User, error)
	FetchUserByID(id uint) (*models.User, error)
	FetchByEmail(email string) (*models.User, error)

	UpdateUser(string, models.UpdateUser) error
	DeleteUser(id uint) error
	VerifyCode(code string) (*models.User, error)
	//UserManagement
	Create(req *models.CreateUser) (*models.User, error)
	AdminUpdate(actorID int64, publicID string, req models.AdminUpdateUser) (*models.User, error)
	AdminDelete(actorID int64, publicID string) error
	FetchUsersPaginated(filter, role, sort string, limit, offset int) ([]models.User, int64, error)
}

type UserServices struct {
	repo       repositories.UserRepository
	mail       *utils.SMTPService
	codeGen    *code.DefaultCodeGenerator
	updateUser models.UpdateUser
}

func NewUserService(repo repositories.UserRepository, mail *utils.SMTPService, codeGen *code.DefaultCodeGenerator, updateUser models.UpdateUser) UserService {
	return &UserServices{repo, mail, codeGen, updateUser}
}

func (s *UserServices) Register(user *models.User) error {
	existingUser, err := s.repo.FetchByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return utils.ErrConflict("email already registered")
	}
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Verified = false
	user.VerificationToken = nil
	user.VerifiedAt = nil
	user.Password = hashed
	user.Role = "researcher"
	user.PublicID = uuid.New()
	if err := s.repo.Create(user); err != nil {
		return err
	}

	rawCode, err := s.codeGen.Generate()
	if err != nil {
		return err
	}
	hashedCode := s.codeGen.Hash(rawCode)
	tokenExpiration := time.Now().Add(5 * time.Hour)

	if err := s.repo.StoreToken(user.Email, hashedCode, tokenExpiration); err != nil {
		return err
	}
	// Link way example
	// link := "http://localhost:3000/api/v1/auth/verify?token=" + token
	// mailBody := "<h1> Click the link to Verify your email</h1>" < "a href=' "+link+"'>click here </a>"

	if err := s.mail.SendVerificationCode(user.Email, rawCode); err != nil {
		// Roll back so the address can register again once email delivery works.
		if delErr := s.repo.UserDelete(uint(user.InternalID)); delErr != nil {
			return fmt.Errorf("send verification email: %w (cleanup failed: %v)", err, delErr)
		}
		return fmt.Errorf("send verification email: %w", err)
	}
	return nil
}

func (s *UserServices) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FetchByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *UserServices) FetchUserByPublicID(publicID string) (*models.User, error) {
	return s.repo.FetchByPublicID(publicID)
}
func (s *UserServices) FetchUserByID(id uint) (*models.User, error) {
	return s.repo.FetchByID(id)
}

// func (s *UserServices) UpdateUser(user *models.User) error {
// 	return s.repo.UserUpdate(user)
// }

func (s *UserServices) UpdateUser(publicID string, req models.UpdateUser) error {
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = &req.Name
	}
	if req.Institution != nil {
		updates["institution"] = &req.Institution
	}
	if req.Institution_Address != nil {
		updates["institution_address"] = &req.Institution_Address
	}
	if req.Position != nil {
		updates["position"] = &req.Position
	}
	// if req.Phone_Number != nil {
	// 	updates["phone_number"] = &req.Phone_Number
	// }
	return s.repo.UserUpdate(publicID, updates)
}
func (s *UserServices) DeleteUser(id uint) error {
	return s.repo.UserDelete(id)
}
func (r *UserServices) FetchByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user)
	return &user, err.Error
}
func (s *UserServices) VerifyCode(code string) (*models.User, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("code is required")
	}

	hashedCode := s.codeGen.Hash(strings.TrimSpace(code))

	var user models.User
	err := config.DB.
		Where("verification_token = ? AND verification_expires_at > ?",
			hashedCode, time.Now()).
		First(&user).Error
	if err != nil {
		return nil, errors.New("invalid or expired code")
	}

	now := time.Now()
	err = config.DB.Model(&user).Updates(map[string]interface{}{
		"verified":                true,
		"verified_at":             now,
		"verification_token":      nil,
		"verification_expires_at": nil,
	}).Error
	if err != nil {
		return nil, err
	}

	user.Verified = true
	user.VerifiedAt = &now
	user.VerificationToken = nil
	user.VerificationExpiry = nil

	return &user, nil
}

//User Management

var ValidRoles = map[string]bool{
	"admin": true, "reviewer": true, "opd": true, "pegawai-brida": true, "media-brida": true, "researcher": true,
}

var userEmailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[A-Za-z]{2,}$`)

func validName(v string) error {
	if n := utf8.RuneCountInString(v); n < 3 || n > 255 {
		return utils.ErrValidation("name must be 3-255 characters")
	}
	return nil
}

func validPassword(v string) error {
	if len(v) < 8 || len(v) > 72 {
		return utils.ErrValidation("password must be 8-72 characters")
	}
	return nil
}

func validText(field, v string, max int) error {
	if utf8.RuneCountInString(v) > max {
		return utils.ErrValidation(fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return nil
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(*p)
}

func (s *UserServices) Create(req *models.CreateUser) (*models.User, error) {
	name := deref(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	institution, address, position := deref(req.Institution), deref(req.Institution_Address), deref(req.Position)
	switch {
	case validName(name) != nil:
		return nil, validName(name)
	case len(email) > 255 || !userEmailPattern.MatchString(email):
		return nil, utils.ErrValidation("email is not valid")
	case validPassword(req.Password) != nil:
		return nil, validPassword(req.Password)
	case !ValidRoles[req.Role]:
		return nil, utils.ErrValidation("role is not valid")
	case validText("institution", institution, 255) != nil:
		return nil, validText("institution", institution, 255)
	case validText("position", position, 100) != nil:
		return nil, validText("position", position, 100)
	}
	if existing, _ := s.repo.FetchByEmail(email); existing != nil && existing.InternalID != 0 {
		return nil, utils.ErrConflict("email already registered")
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := models.User{
		PublicID: uuid.New(), Name: &name, Email: email, Password: hashed,
		Institution: institution, Institution_Address: address, Position: position,
		Role: req.Role, Verified: true, VerifiedAt: &now,
	}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServices) loadTarget(publicID string) (*models.User, error) {
	user, err := s.repo.FetchByPublicID(publicID)
	if err != nil || user == nil || user.InternalID == 0 {
		return nil, utils.ErrNotFound("user not found")
	}
	return user, nil
}

// ensureAnotherAdmin refuses changes that would leave the system without an administrator.
func (s *UserServices) ensureAnotherAdmin() error {
	n, err := s.repo.CountAdmins()
	if err != nil {
		return err
	}
	if n <= 1 {
		return utils.ErrConflict("at least one administrator must remain")
	}
	return nil
}

func (s *UserServices) AdminUpdate(actorID int64, publicID string, req models.AdminUpdateUser) (*models.User, error) {
	target, err := s.loadTarget(publicID)
	if err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if err := validName(name); err != nil {
			return nil, err
		}
		updates["name"] = name
	}
	if req.Institution != nil {
		v := strings.TrimSpace(*req.Institution)
		if err := validText("institution", v, 255); err != nil {
			return nil, err
		}
		updates["institution"] = v
	}
	if req.Institution_Address != nil {
		updates["institution_address"] = strings.TrimSpace(*req.Institution_Address)
	}
	if req.Position != nil {
		v := strings.TrimSpace(*req.Position)
		if err := validText("position", v, 100); err != nil {
			return nil, err
		}
		updates["position"] = v
	}
	if req.Role != nil && *req.Role != target.Role {
		if !ValidRoles[*req.Role] {
			return nil, utils.ErrValidation("role is not valid")
		}
		if target.InternalID == actorID {
			return nil, utils.ErrForbidden("you cannot change your own role")
		}
		if target.Role == "admin" {
			if err := s.ensureAnotherAdmin(); err != nil {
				return nil, err
			}
		}
		updates["role"] = *req.Role
	}
	if req.Verified != nil && *req.Verified != target.Verified {
		updates["verified"] = *req.Verified
		if *req.Verified {
			updates["verified_at"] = time.Now()
		} else {
			updates["verified_at"] = nil
		}
	}
	if req.Password != nil && *req.Password != "" {
		if err := validPassword(*req.Password); err != nil {
			return nil, err
		}
		hashed, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		updates["password"] = hashed
	}
	if len(updates) == 0 {
		return target, nil
	}
	if err := s.repo.UserUpdate(publicID, updates); err != nil {
		return nil, err
	}
	return s.loadTarget(publicID)
}

func (s *UserServices) AdminDelete(actorID int64, publicID string) error {
	target, err := s.loadTarget(publicID)
	if err != nil {
		return err
	}
	if target.InternalID == actorID {
		return utils.ErrForbidden("you cannot delete your own account")
	}
	if target.Role == "admin" {
		if err := s.ensureAnotherAdmin(); err != nil {
			return err
		}
	}
	n, err := s.repo.UserDependencyCount(target.InternalID)
	if err != nil {
		return err
	}
	if n > 0 {
		return utils.ErrConflict("user still owns research, news, documentation or review data and cannot be deleted")
	}
	return s.repo.UserDelete(uint(target.InternalID))
}

func (s *UserServices) FetchUsersPaginated(filter, role, sort string, limit, offset int) ([]models.User, int64, error) {
	return s.repo.FetchUsersPaginated(filter, role, sort, limit, offset)
}
