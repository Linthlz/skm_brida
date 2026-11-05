package controllers

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[A-Za-z]{2,}$`)

type registerRequest struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	PhoneNumber        string `json:"phone_number"`
	Institution        string `json:"institution"`
	InstitutionAddress string `json:"institution_address"`
	Position           string `json:"position"`
}

func (r *registerRequest) validate() string {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
	r.Institution = strings.TrimSpace(r.Institution)
	r.InstitutionAddress = strings.TrimSpace(r.InstitutionAddress)
	r.Position = strings.TrimSpace(r.Position)
	switch {
	case utf8.RuneCountInString(r.Name) < 3 || utf8.RuneCountInString(r.Name) > 255:
		return "name must be 3-255 characters"
	case len(r.Email) > 255 || !emailPattern.MatchString(r.Email):
		return "email is not valid"
	case len(r.Password) < 8 || len(r.Password) > 72:
		return "password must be 8-72 characters"
	case len(r.PhoneNumber) > 50:
		return "phone_number must be at most 50 characters"
	case utf8.RuneCountInString(r.Position) > 100:
		return "position must be at most 100 characters"
	}
	return ""
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	if !repositories.SettingBool("registration_open") {
		return utils.Forbidden(ctx, "Registration closed", "self registration is currently disabled by the administrator")
	}
	var body registerRequest
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	if msg := body.validate(); msg != "" {
		return utils.BadReq(ctx, "Registration failed", msg)
	}
	user := &models.User{
		Name:                &body.Name,
		Email:               body.Email,
		Password:            body.Password,
		Phone_Number:        body.PhoneNumber,
		Institution:         body.Institution,
		Institution_Address: body.InstitutionAddress,
		Position:            body.Position,
	}
	if err := c.service.Register(user); err != nil {
		return utils.HandleError(ctx, "Registration failed", err)
	}
	var UserResponse models.UserResponse
	_ = copier.Copy(&UserResponse, user)
	return utils.CreatedWithData(ctx, "Registration success, check your email for the verification code", UserResponse)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	body.Email = strings.ToLower(strings.TrimSpace(body.Email))
	if body.Email == "" || body.Password == "" {
		return utils.BadReq(ctx, "Login failed", "email and password are required")
	}
	user, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.UnauthorizedAccess(ctx, "Login Failed", err.Error())
	}
	var UserResponse models.UserResponse
	_ = copier.Copy(&UserResponse, user)
	Token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.Verified, user.PublicID)
	Refresh_Token, _ := utils.RefreshToken(user.InternalID)
	return utils.Success(ctx, "Login success!", fiber.Map{
		"Token":         Token,
		"refresh_token": Refresh_Token,
		"Data":          UserResponse,
	})
}

// Fix
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadReq(ctx, "Invalid ID Format", err.Error())
	}
	var UpdateUser models.UpdateUser
	if err := ctx.BodyParser(&UpdateUser); err != nil {
		return utils.BadReq(ctx, "Failed Request Parsing", err.Error())
	}
	// if err := ctx.BodyParser(&UpdateUser); err != nil {
	// 	return utils.InternalServerError(ctx, "Failed to Update user Data", err.Error())
	// }
	// user, err := c.service.FetchUserByPublicID(publicID.String())
	if err := c.service.UpdateUser(publicID.String(), UpdateUser); err != nil {
		return utils.InternalServerError(ctx, "Update Failed", err.Error())
	}
	return utils.Success(ctx, "User Updated", &UpdateUser)
	// var UserResponse models.User
	// err = copier.Copy(&UserResponse, userUpdated)
	// if err != nil {
	// 	return utils.InternalServerError(ctx, "error data parsing", err.Error())
	// }
}

func (c *UserController) FetchProfileData(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.FetchUserByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "User not found", err.Error())
	}
	var UserResponse models.UserResponse
	err = copier.Copy(&UserResponse, user)
	if err != nil {
		return utils.InternalServerError(ctx, "Internal server error", err.Error())
	}
	return utils.Success(ctx, "User Data Fetched", UserResponse)
}

// ================= User Management ============================
func toUserResponse(u *models.User) models.UserResponse {
	r := models.UserResponse{
		PublicID: u.PublicID, Name: u.Name, Email: u.Email, Institution: u.Institution,
		Institution_Address: u.Institution_Address, Position: u.Position, Role: u.Role,
		Verified: u.Verified, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
	if u.VerifiedAt != nil {
		r.VerifiedAt = *u.VerifiedAt
	}
	return r
}

// CreateUser: POST /api/v1/admin/user
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req models.CreateUser
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	user, err := c.service.Create(&req)
	if err != nil {
		return utils.HandleError(ctx, "Failed to create user", err)
	}
	return utils.CreatedWithData(ctx, "User created", toUserResponse(user))
}

// UpdateUserByAdmin: PATCH /api/v1/admin/user/:id
func (c *UserController) UpdateUserByAdmin(ctx *fiber.Ctx) error {
	actorID, ok := currentUserID(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Unauthorized", "invalid token")
	}
	id := ctx.Params("id")
	if !isUUID(id) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	var req models.AdminUpdateUser
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	user, err := c.service.AdminUpdate(actorID, id, req)
	if err != nil {
		return utils.HandleError(ctx, "Failed to update user", err)
	}
	return utils.Success(ctx, "User updated", toUserResponse(user))
}

// DeleteUserByAdmin: DELETE /api/v1/admin/user/:id
func (c *UserController) DeleteUserByAdmin(ctx *fiber.Ctx) error {
	actorID, ok := currentUserID(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Unauthorized", "invalid token")
	}
	id := ctx.Params("id")
	if !isUUID(id) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	if err := c.service.AdminDelete(actorID, id); err != nil {
		return utils.HandleError(ctx, "Failed to delete user", err)
	}
	return utils.Success(ctx, "User deleted", nil)
}

// Users: GET /api/v1/admin/user?page=&limit=&filter=&role=
func (c *UserController) Users(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 20)
	filter := strings.TrimSpace(ctx.Query("filter", ""))
	role := strings.TrimSpace(ctx.Query("role", ""))
	sort := ctx.Query("sort", "")
	if len(filter) > 100 {
		return utils.BadReq(ctx, "Invalid filter", "filter is too long")
	}
	users, total, err := c.service.FetchUsersPaginated(filter, role, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch Users", err.Error())
	}
	response := make([]models.UserResponse, 0, len(users))
	for i := range users {
		response = append(response, toUserResponse(&users[i]))
	}
	return utils.SuccessPaginated(ctx, "User Data Fetched", response, metaFor(page, limit, total, filter, sort))
}
