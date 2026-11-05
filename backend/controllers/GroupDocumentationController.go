package controllers

import (
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type GroupDocumentationController struct {
	service services.GroupDocumentationService
}

func NewGroupDocumentationController(service services.GroupDocumentationService) *GroupDocumentationController {
	return &GroupDocumentationController{service: service}
}

func (c *GroupDocumentationController) CreateDocumentation(ctx *fiber.Ctx) error {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	FuserID, ok := claims["user_id"].(float64)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	userID := int64(FuserID)
	groupPublicID := ctx.Params("id")
	if groupPublicID == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}
	var input models.GroupDocumentation
	input.Title = ctx.FormValue("title")
	input.Narrative = ctx.FormValue("narrative")
	input.Location = ctx.FormValue("location")
	if v := ctx.FormValue("video_url"); v != "" {
		input.VideoURL = &v
	}
	if d := ctx.FormValue("event_date"); d != "" {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			return utils.BadReq(ctx, "Invalid event_date", "Use YYYY-MM-DD format")
		}
		input.EventDate = &t
	}
	form, err := ctx.MultipartForm()
	if err != nil {
		return utils.BadReq(ctx, "Failed to read form", err.Error())
	}
	if len(form.File["images"]) == 0 {
		return utils.BadReq(ctx, "Images required", "At least one image is required")
	}
	if len(form.File["images"]) > utils.MaxImagesPerDocumentation {
		return utils.BadReq(ctx, "Too many images", "Maximum 10 images per documentation")
	}
	imagePaths, err := utils.SaveImageBatch(form, "images", "./public/storage/documentations", "/storage/documentations")
	if err != nil {
		return utils.BadReq(ctx, "Image upload failed", err.Error())
	}
	captions := form.Value["captions"]
	doc, err := c.service.Create(userID, groupPublicID, input, imagePaths, captions)
	if err != nil {
		return utils.BadReq(ctx, "Cannot create documentation", err.Error())
	}
	return utils.Success(ctx, "Documentation created", doc)
}

func (c *GroupDocumentationController) GetDocumentation(ctx *fiber.Ctx) error {
	docPublicID := ctx.Params("docID")
	if docPublicID == "" {
		return utils.NotFound(ctx, "Documentation Not Found", "")
	}
	doc, err := c.service.FetchByPublicID(docPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Documentation Not Found", err.Error())
	}
	return utils.Success(ctx, "Documentation fetched", doc)
}

func (c *GroupDocumentationController) GetGroupDocumentations(ctx *fiber.Ctx) error {
	groupPublicID := ctx.Params("id")
	if groupPublicID == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	docs, total, err := c.service.FetchByGroup(groupPublicID, offset, limit)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch", err.Error())
	}
	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
	}
	return utils.SuccessPaginated(ctx, "Documentations fetched", docs, meta)
}

func (c *GroupDocumentationController) UpdateDocumentation(ctx *fiber.Ctx) error {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	FuserID, ok := claims["user_id"].(float64)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	userID := int64(FuserID)
	docPublicID := ctx.Params("docID")
	if docPublicID == "" {
		return utils.NotFound(ctx, "Documentation Not Found", "")
	}
	var body struct {
		Title     string `json:"title"`
		Narrative string `json:"narrative"`
		Location  string `json:"location"`
		VideoURL  string `json:"video_url"`
		EventDate string `json:"event_date"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	input := models.GroupDocumentation{
		Title:     body.Title,
		Narrative: body.Narrative,
		Location:  body.Location,
	}
	if body.VideoURL != "" {
		input.VideoURL = &body.VideoURL
	}
	if body.EventDate != "" {
		t, err := time.Parse("2006-01-02", body.EventDate)
		if err != nil {
			return utils.BadReq(ctx, "Invalid event_date", "Use YYYY-MM-DD format")
		}
		input.EventDate = &t
	}
	if err := c.service.Update(userID, docPublicID, input); err != nil {
		return utils.BadReq(ctx, "Cannot update documentation", err.Error())
	}
	return utils.Success(ctx, "Documentation updated", nil)
}

func (c *GroupDocumentationController) DeleteDocumentation(ctx *fiber.Ctx) error {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	FuserID, ok := claims["user_id"].(float64)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	userID := int64(FuserID)
	docPublicID := ctx.Params("docID")
	if docPublicID == "" {
		return utils.NotFound(ctx, "Documentation Not Found", "")
	}
	if err := c.service.Delete(userID, docPublicID); err != nil {
		return utils.BadReq(ctx, "Cannot delete documentation", err.Error())
	}
	return utils.Success(ctx, "Documentation deleted", nil)
}
