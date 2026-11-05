package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsXanadium/singa-riset-backend/application"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type VerificationHandler struct {
	verificationService *application.VerificationService
	user                services.UserService
}

func NewVerificationHandler(service *application.VerificationService, userService services.UserService) *VerificationHandler {
	return &VerificationHandler{verificationService: service, user: userService}
}

func (h *VerificationHandler) SendVerification(c *fiber.Ctx) error {
	email := c.Params("email")
	if err := h.verificationService.SendVerification(email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"{+}": "verification code sent"})
}

//	func (h *VerificationHandler) CheckCode(c *fiber.Ctx) error {
//		// email := c.Params("email")
//		// code := c.Params("code")
//		var body struct {
//			Email string `json:"email"`
//			Code  string `json:"code"`
//		}
//		if err := c.BodyParser(&body); err != nil {
//			return utils.BadReq(c, "Invalid Request", err.Error())
//		}
//		if err := h.verificationService.VerifyCode(body.Email, body.Code); err != nil {
//			return utils.BadReq(c, "Invalid Request", err.Error())
//		}
//		user, err := h.user.FetchByEmail(body.Email)
//		if err != nil {
//			return utils.BadReq(c, "Invalid Request", err.Error())
//		}
//		if user == nil {
//			return utils.NotFound(c, "Error:", "user doesnt exist")
//		}
//		Token, err := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.Verified, user.PublicID)
//		if err != nil {
//			return utils.BadReq(c, "Invalid Request", err.Error())
//		}
//		return c.JSON(fiber.Map{
//			"message": "Email Verified!",
//			"Token":   Token})
//	}
func (h *VerificationHandler) CheckCode(c *fiber.Ctx) error {
	var body struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.BadReq(c, "Invalid Request", err.Error())
	}
	if body.Code == "" {
		return utils.BadReq(c, "Invalid Request", "code is required")
	}

	user, err := h.user.VerifyCode(body.Code)
	if err != nil || user == nil {
		return utils.BadReq(c, "Verification failed", "invalid or expired code")
	}

	token, err := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.Verified, user.PublicID)
	if err != nil {
		return utils.InternalServerError(c, "Verification failed", err.Error())
	}
	return utils.Success(c, "Email Verified!", fiber.Map{"token": token})
}
