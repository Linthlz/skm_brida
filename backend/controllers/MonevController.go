package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type MonevController struct {
	service services.MonevService
}

func NewMonevController(service services.MonevService) *MonevController {
	return &MonevController{service: service}
}

func (c *MonevController) getUserIDAndRole(ctx *fiber.Ctx) (int64, string, bool) {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return 0, "", false
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", false
	}
	FuserID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", false
	}
	role, _ := claims["role"].(string)
	return int64(FuserID), role, true
}

func (c *MonevController) CreateBatch(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	var body struct {
		Year        int    `json:"year"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	b, err := c.service.CreateBatch(userID, models.MonevBatch{
		Year: body.Year, Title: body.Title, Description: body.Description, Status: body.Status,
	})
	if err != nil {
		return utils.BadReq(ctx, "Cannot create batch", err.Error())
	}
	return utils.Success(ctx, "Batch created", b)
}

func (c *MonevController) ListBatches(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")
	list, total, err := c.service.ListBatches(filter, sort, offset, limit)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch", err.Error())
	}
	meta := utils.PaginationMeta{
		Page: page, Limit: limit, Total: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter, Sort: sort,
	}
	return utils.SuccessPaginated(ctx, "Batches fetched", list, meta)
}

func (c *MonevController) GetBatch(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if publicID == "" {
		return utils.NotFound(ctx, "Batch Not Found", "")
	}
	b, err := c.service.FetchBatchByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Batch Not Found", err.Error())
	}
	return utils.Success(ctx, "Batch fetched", b)
}

func (c *MonevController) UpdateBatch(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	publicID := ctx.Params("id")
	var body struct {
		Year        int    `json:"year"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	if err := c.service.UpdateBatch(userID, publicID, models.MonevBatch{
		Year: body.Year, Title: body.Title, Description: body.Description, Status: body.Status,
	}); err != nil {
		return utils.BadReq(ctx, "Cannot update batch", err.Error())
	}
	return utils.Success(ctx, "Batch updated", nil)
}

func (c *MonevController) DeleteBatch(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	publicID := ctx.Params("id")
	if err := c.service.DeleteBatch(userID, publicID); err != nil {
		return utils.BadReq(ctx, "Cannot delete batch", err.Error())
	}
	return utils.Success(ctx, "Batch deleted", nil)
}

func (c *MonevController) AddOpd(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	batchPublicID := ctx.Params("id")
	var body struct {
		OpdPublicID string `json:"opd_public_id"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	entry, err := c.service.AddOpdToBatch(userID, batchPublicID, body.OpdPublicID)
	if err != nil {
		return utils.BadReq(ctx, "Cannot add OPD", err.Error())
	}
	return utils.Success(ctx, "OPD added to batch", entry)
}

func (c *MonevController) ListBatchEntries(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	batchPublicID := ctx.Params("id")
	entries, err := c.service.ListBatchEntries(userID, batchPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Not Found", err.Error())
	}
	return utils.Success(ctx, "Entries fetched", entries)
}

func (c *MonevController) RemoveEntry(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	entryPublicID := ctx.Params("entryID")
	if err := c.service.RemoveOpdFromBatch(userID, entryPublicID); err != nil {
		return utils.BadReq(ctx, "Cannot remove", err.Error())
	}
	return utils.Success(ctx, "Entry removed", nil)
}

func (c *MonevController) GetEntry(ctx *fiber.Ctx) error {
	userID, role, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	entryPublicID := ctx.Params("entryID")
	e, err := c.service.FetchEntryByPublicID(userID, role, entryPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Entry Not Found", err.Error())
	}
	return utils.Success(ctx, "Entry fetched", e)
}

func (c *MonevController) ListMyEntries(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	// Ambil opd ID dari user
	user, err := c.service.GetUserOpd(userID)
	if err != nil || user == nil {
		return utils.NotFound(ctx, "User not found", "")
	}
	if user.OpdInternalID == nil {
		return utils.BadReq(ctx, "No OPD assigned", "Your account has no OPD assigned")
	}
	entries, err := c.service.ListMyEntries(userID, *user.OpdInternalID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch", err.Error())
	}
	return utils.Success(ctx, "Entries fetched", entries)
}

func (c *MonevController) FillRepresentative(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	entryPublicID := ctx.Params("entryID")
	var body struct {
		Name string `json:"name"`
		Nip  string `json:"nip"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	if err := c.service.FillRepresentative(userID, entryPublicID, body.Name, body.Nip); err != nil {
		return utils.BadReq(ctx, "Cannot fill representative", err.Error())
	}
	return utils.Success(ctx, "Representative saved", nil)
}

func (c *MonevController) SubmitEntry(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	entryPublicID := ctx.Params("entryID")
	if err := c.service.SubmitEntry(userID, entryPublicID); err != nil {
		return utils.BadReq(ctx, "Cannot submit", err.Error())
	}
	return utils.Success(ctx, "Entry submitted", nil)
}

func (c *MonevController) VerifyEntry(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	entryPublicID := ctx.Params("entryID")
	if err := c.service.VerifyEntry(userID, entryPublicID); err != nil {
		return utils.BadReq(ctx, "Cannot verify", err.Error())
	}
	return utils.Success(ctx, "Entry verified", nil)
}

func (c *MonevController) AddKajian(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	entryPublicID := ctx.Params("entryID")
	var body struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		DisplayOrder int    `json:"display_order"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	k, err := c.service.AddKajian(userID, entryPublicID, models.KajianResearch{
		Title: body.Title, Description: body.Description, DisplayOrder: body.DisplayOrder,
	})
	if err != nil {
		return utils.BadReq(ctx, "Cannot add kajian", err.Error())
	}
	return utils.Success(ctx, "Kajian added", k)
}

func (c *MonevController) UpdateKajian(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	kajianPublicID := ctx.Params("kajianID")
	var body struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		DisplayOrder int    `json:"display_order"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	if err := c.service.UpdateKajian(userID, kajianPublicID, models.KajianResearch{
		Title: body.Title, Description: body.Description, DisplayOrder: body.DisplayOrder,
	}); err != nil {
		return utils.BadReq(ctx, "Cannot update kajian", err.Error())
	}
	return utils.Success(ctx, "Kajian updated", nil)
}

func (c *MonevController) DeleteKajian(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	kajianPublicID := ctx.Params("kajianID")
	if err := c.service.DeleteKajian(userID, kajianPublicID); err != nil {
		return utils.BadReq(ctx, "Cannot delete kajian", err.Error())
	}
	return utils.Success(ctx, "Kajian deleted", nil)
}

func (c *MonevController) AddPoint(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	kajianPublicID := ctx.Params("kajianID")
	var body struct {
		Title string `json:"title"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	p, err := c.service.AddPoint(userID, kajianPublicID, models.ResearchPoint{Title: body.Title})
	if err != nil {
		return utils.BadReq(ctx, "Cannot add point", err.Error())
	}
	return utils.Success(ctx, "Point added", p)
}

func (c *MonevController) UpdatePoint(ctx *fiber.Ctx) error {
	userID, role, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	pointPublicID := ctx.Params("pointID")
	var body struct {
		Title                         string  `json:"title"`
		MonitoringFilePath            *string `json:"monitoring_file_path"`
		HasMonitoringResult           bool    `json:"has_monitoring_result"`
		MonitoringResultDescription   string  `json:"monitoring_result_description"`
		KendalaPelaksanaanRekomendasi string  `json:"kendala_pelaksanaan_rekomendasi"`
		Benefits                      string  `json:"benefits"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	input := models.ResearchPoint{
		Title:                         body.Title,
		MonitoringFilePath:            body.MonitoringFilePath,
		HasMonitoringResult:           body.HasMonitoringResult,
		MonitoringResultDescription:   body.MonitoringResultDescription,
		KendalaPelaksanaanRekomendasi: body.KendalaPelaksanaanRekomendasi,
		Benefits:                      body.Benefits,
	}
	if err := c.service.UpdatePoint(userID, role, pointPublicID, input); err != nil {
		return utils.BadReq(ctx, "Cannot update point", err.Error())
	}
	return utils.Success(ctx, "Point updated", nil)
}

func (c *MonevController) DeletePoint(ctx *fiber.Ctx) error {
	userID, _, ok := c.getUserIDAndRole(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	pointPublicID := ctx.Params("pointID")
	if err := c.service.DeletePoint(userID, pointPublicID); err != nil {
		return utils.BadReq(ctx, "Cannot delete point", err.Error())
	}
	return utils.Success(ctx, "Point deleted", nil)
}
