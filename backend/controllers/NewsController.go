package controllers

import (
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

const (
	newsCoverDiskDir = "./public/storage/news"
	newsCoverURLDir  = "/storage/news"
	maxFilterLength  = 100
)

type NewsController struct {
	service services.NewsService
}

func NewNewsController(service services.NewsService) *NewsController {
	return &NewsController{service: service}
}

type newsForm struct {
	Title            string  `json:"title" form:"title"`
	CategoryPublicID string  `json:"category_public_id" form:"category_public_id"`
	Content          string  `json:"content" form:"content"`
	Source           *string `json:"source" form:"source"`
	Status           string  `json:"status" form:"status"`
}

func isMultipart(ctx *fiber.Ctx) bool {
	return strings.HasPrefix(string(ctx.Request().Header.ContentType()), fiber.MIMEMultipartForm)
}

// saveCover stores an optional single "cover_image" upload and returns its public URL path.
func saveCover(ctx *fiber.Ctx, required bool) (string, error) {
	if !isMultipart(ctx) {
		if required {
			return "", utils.ErrValidation("request must be multipart/form-data with a cover_image file")
		}
		return "", nil
	}
	form, err := ctx.MultipartForm()
	if err != nil {
		return "", utils.ErrValidation("failed to read form data")
	}
	files := form.File["cover_image"]
	if len(files) == 0 {
		if required {
			return "", utils.ErrValidation("cover_image is required")
		}
		return "", nil
	}
	if len(files) > 1 {
		return "", utils.ErrValidation("only one cover_image is allowed")
	}
	paths, err := utils.SaveImageBatch(form, "cover_image", newsCoverDiskDir, newsCoverURLDir)
	if err != nil {
		return "", utils.ErrValidation(err.Error())
	}
	return paths[0], nil
}

func removeCover(urlPath string) {
	if !strings.HasPrefix(urlPath, newsCoverURLDir+"/") {
		return
	}
	if err := utils.FileDeletion("./public" + urlPath); err != nil {
		log.Printf("[news] failed to remove cover %s: %v", urlPath, err)
	}
}

func parsePaging(ctx *fiber.Ctx) (page, limit, offset int) {
	page, _ = strconv.Atoi(ctx.Query("page", "1"))
	limit, _ = strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit, (page - 1) * limit
}

func (c *NewsController) list(ctx *fiber.Ctx, publishedOnly bool) error {
	page, limit, offset := parsePaging(ctx)
	filter := strings.TrimSpace(ctx.Query("filter", ""))
	if len(filter) > maxFilterLength {
		return utils.BadReq(ctx, "Invalid filter", "filter is too long")
	}
	categorySlug := strings.TrimSpace(ctx.Query("category", ""))
	sort := ctx.Query("sort", "")
	if sort != "" && sort != "oldest" && sort != "title" {
		return utils.BadReq(ctx, "Invalid sort", "sort must be 'oldest' or 'title'")
	}
	news, total, err := c.service.List(filter, categorySlug, sort, offset, limit, publishedOnly)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch news", err.Error())
	}
	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}
	return utils.SuccessPaginated(ctx, "News fetched", news, meta)
}

// List is the public listing: published news only.
func (c *NewsController) List(ctx *fiber.Ctx) error {
	return c.list(ctx, true)
}

// ListAll is the editor listing: drafts included unless published_only=true.
func (c *NewsController) ListAll(ctx *fiber.Ctx) error {
	return c.list(ctx, ctx.Query("published_only", "false") == "true")
}

func (c *NewsController) GetBySlug(ctx *fiber.Ctx) error {
	slug := strings.TrimSpace(ctx.Params("slug"))
	if slug == "" {
		return utils.BadReq(ctx, "Invalid slug", "slug is required")
	}
	news, err := c.service.FetchBySlug(slug, true)
	if err != nil {
		return utils.HandleError(ctx, "News not found", err)
	}
	return utils.Success(ctx, "News fetched", news)
}

func (c *NewsController) GetByID(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if !isUUID(publicID) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	news, err := c.service.FetchByPublicID(publicID)
	if err != nil {
		return utils.HandleError(ctx, "News not found", err)
	}
	return utils.Success(ctx, "News fetched", news)
}

func (c *NewsController) CreateNews(ctx *fiber.Ctx) error {
	userID, ok := currentUserID(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Unauthorized", "invalid token")
	}
	var body newsForm
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse form data")
	}
	cover, err := saveCover(ctx, true)
	if err != nil {
		return utils.HandleError(ctx, "Cannot create news", err)
	}
	news, err := c.service.Create(userID, body.CategoryPublicID, models.News{
		Title:      body.Title,
		Content:    body.Content,
		CoverImage: cover,
		Source:     body.Source,
		Status:     body.Status,
	})
	if err != nil {
		removeCover(cover)
		return utils.HandleError(ctx, "Cannot create news", err)
	}
	return utils.CreatedWithData(ctx, "News created", news)
}

// Update accepts JSON, or multipart/form-data when a new cover_image is uploaded.
func (c *NewsController) Update(ctx *fiber.Ctx) error {
	userID, ok := currentUserID(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Unauthorized", "invalid token")
	}
	publicID := ctx.Params("id")
	if !isUUID(publicID) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	existing, err := c.service.FetchByPublicID(publicID)
	if err != nil {
		return utils.HandleError(ctx, "Cannot update news", err)
	}
	var body newsForm
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	cover, err := saveCover(ctx, false)
	if err != nil {
		return utils.HandleError(ctx, "Cannot update news", err)
	}
	news, err := c.service.Update(userID, publicID, body.CategoryPublicID, models.News{
		Title:      body.Title,
		Content:    body.Content,
		CoverImage: cover,
		Source:     body.Source,
		Status:     body.Status,
	})
	if err != nil {
		removeCover(cover)
		return utils.HandleError(ctx, "Cannot update news", err)
	}
	if cover != "" && cover != existing.CoverImage {
		removeCover(existing.CoverImage)
	}
	return utils.Success(ctx, "News updated", news)
}

func (c *NewsController) Delete(ctx *fiber.Ctx) error {
	userID, ok := currentUserID(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Unauthorized", "invalid token")
	}
	publicID := ctx.Params("id")
	if !isUUID(publicID) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	if err := c.service.Delete(userID, publicID); err != nil {
		return utils.HandleError(ctx, "Cannot delete news", err)
	}
	return utils.Success(ctx, "News deleted", nil)
}
