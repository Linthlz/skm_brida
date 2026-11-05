package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type ResearchDeliverableController struct {
	service services.ResearchDeliverableService
}

func NewResearchDeliverableController(service services.ResearchDeliverableService) *ResearchDeliverableController {
	return &ResearchDeliverableController{service: service}
}

func (c *ResearchDeliverableController) getUserID(ctx *fiber.Ctx) (int64, error) {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return 0, fiber.ErrUnauthorized
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	FuserID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return int64(FuserID), nil
}

func (c *ResearchDeliverableController) List(ctx *fiber.Ctx) error {
	groupPublicID := ctx.Params("id")
	if groupPublicID == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}
	list, progress, err := c.service.FetchByGroup(groupPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Not Found", err.Error())
	}
	return utils.Success(ctx, "Deliverables fetched", fiber.Map{
		"progress":     progress,
		"deliverables": list,
	})
}

func (c *ResearchDeliverableController) Create(ctx *fiber.Ctx) error {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	groupPublicID := ctx.Params("id")
	if groupPublicID == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ProofType   string `json:"proof_type"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	input := models.ResearchDeliverable{
		Title:       body.Title,
		Description: body.Description,
		ProofType:   body.ProofType,
	}
	d, err := c.service.Create(userID, groupPublicID, input, nil)
	if err != nil {
		return utils.BadReq(ctx, "Cannot create deliverable", err.Error())
	}
	return utils.Success(ctx, "Deliverable created", d)
}

func (c *ResearchDeliverableController) Update(ctx *fiber.Ctx) error {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	publicID := ctx.Params("dID")
	if publicID == "" {
		return utils.NotFound(ctx, "Not Found", "")
	}
	var body struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		ProofType    string `json:"proof_type"`
		DisplayOrder int    `json:"display_order"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	input := models.ResearchDeliverable{
		Title:        body.Title,
		Description:  body.Description,
		ProofType:    body.ProofType,
		DisplayOrder: body.DisplayOrder,
	}
	if err := c.service.Update(userID, publicID, input); err != nil {
		return utils.BadReq(ctx, "Cannot update", err.Error())
	}
	return utils.Success(ctx, "Deliverable updated", nil)
}

func (c *ResearchDeliverableController) MarkDone(ctx *fiber.Ctx) error {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	publicID := ctx.Params("dID")
	if publicID == "" {
		return utils.NotFound(ctx, "Not Found", "")
	}
	form, formErr := ctx.MultipartForm()
	var filePath *string

	if formErr == nil && len(form.File["file"]) > 0 {
		paths, saveErr := utils.SaveFileBatch(form, "file", "./public/storage/deliverables", "/storage/deliverables")
		if saveErr != nil {
			return utils.BadReq(ctx, "File upload failed", saveErr.Error())
		}
		if len(paths) > 0 {
			filePath = &paths[0]
		}
	}

	if err := c.service.MarkDone(userID, publicID, filePath); err != nil {
		return utils.BadReq(ctx, "Cannot mark done", err.Error())
	}
	return utils.Success(ctx, "Deliverable marked done", nil)
}

func (c *ResearchDeliverableController) MarkPending(ctx *fiber.Ctx) error {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	publicID := ctx.Params("dID")
	if publicID == "" {
		return utils.NotFound(ctx, "Not Found", "")
	}
	if err := c.service.MarkPending(userID, publicID); err != nil {
		return utils.BadReq(ctx, "Cannot mark pending", err.Error())
	}
	return utils.Success(ctx, "Deliverable marked pending", nil)
}

func (c *ResearchDeliverableController) Delete(ctx *fiber.Ctx) error {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	publicID := ctx.Params("dID")
	if publicID == "" {
		return utils.NotFound(ctx, "Not Found", "")
	}
	if err := c.service.Delete(userID, publicID); err != nil {
		return utils.BadReq(ctx, "Cannot delete", err.Error())
	}
	return utils.Success(ctx, "Deliverable deleted", nil)
}
