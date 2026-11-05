package controllers

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

// Research status are {pending, under-review, accepted, revision, rejected}
type ResearchController struct {
	services services.ResearchService
	GService services.ResearchGroupService
}

func NewResearchController(s services.ResearchService, GS services.ResearchGroupService) *ResearchController {
	return &ResearchController{services: s, GService: GS}
}

// researchSubmitRequest whitelists the fields a researcher may send; status,
// reviewer and codes are always set by the server.
type researchSubmitRequest struct {
	Title                        string `form:"title" json:"title"`
	Purpose                      string `form:"purpose" json:"purpose"`
	Category                     string `form:"category" json:"category"`
	FundingScheme                string `form:"funding_scheme" json:"funding_scheme"`
	FundAmount                   string `form:"fund_amount" json:"fund_amount"`
	Location                     string `form:"location" json:"location"`
	Address                      string `form:"address" json:"address"`
	RPJMD                        string `form:"rpjmd" json:"rpjmd"`
	PromisedOutput               string `form:"promised_output" json:"promised_output"`
	ResearchSignificance         string `form:"research_significance" json:"research_significance"`
	ResearchTimeRange            string `form:"research_time_range" json:"research_time_range"`
	PartnershipTargetSelection   string `form:"partnership_target_selection" json:"partnership_target_selection"`
	PartnerName                  string `form:"partner_name" json:"partner_name"`
	EstimatedDirectBeneficiaries string `form:"estimated_number_of_direct_beneficiaries" json:"estimated_number_of_direct_beneficiaries"`
}

func (r *researchSubmitRequest) validate() string {
	trim := func(s *string) { *s = strings.TrimSpace(*s) }
	for _, f := range []*string{&r.Title, &r.Purpose, &r.Category, &r.FundingScheme, &r.FundAmount, &r.Location, &r.Address,
		&r.RPJMD, &r.PromisedOutput, &r.ResearchSignificance, &r.ResearchTimeRange, &r.PartnershipTargetSelection,
		&r.PartnerName, &r.EstimatedDirectBeneficiaries} {
		trim(f)
	}
	required := map[string]string{
		"title": r.Title, "purpose": r.Purpose, "category": r.Category, "funding_scheme": r.FundingScheme,
		"fund_amount": r.FundAmount, "location": r.Location, "rpjmd": r.RPJMD, "promised_output": r.PromisedOutput,
	}
	for _, key := range []string{"title", "purpose", "category", "funding_scheme", "fund_amount", "location", "rpjmd", "promised_output"} {
		if required[key] == "" {
			return key + " is required"
		}
	}
	if utf8.RuneCountInString(r.Title) > 255 {
		return "title must be at most 255 characters"
	}
	for _, s := range []string{r.Category, r.FundingScheme, r.FundAmount, r.Location, r.ResearchTimeRange, r.PartnerName, r.EstimatedDirectBeneficiaries} {
		if utf8.RuneCountInString(s) > 255 {
			return "text fields must be at most 255 characters"
		}
	}
	if _, err := strconv.ParseUint(r.FundAmount, 10, 64); err != nil {
		return "fund_amount must be a whole number in rupiah"
	}
	return ""
}

func (c *ResearchController) ResearchSubmission(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	pubID, ok := claims["pub_id"].(string)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Invalid Token", "no pub id")
	}
	if !repositories.SettingBool("batch_open") {
		return utils.Forbidden(ctx, "Submission batch closed", "the submission batch is currently closed")
	}
	var body researchSubmitRequest
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse form data")
	}
	if msg := body.validate(); msg != "" {
		return utils.BadReq(ctx, "Failed to Submit Research Proposal", msg)
	}
	filename, err := utils.FileUpload(ctx, "research_proposal_file", "./public/storage/research/proposal")
	if err != nil {
		return utils.BadReq(ctx, "Unable to Upload File", err.Error())
	}
	if filename == "" {
		return utils.BadReq(ctx, "Failed to Submit Research Proposal", "research_proposal_file (PDF) is required")
	}
	title := body.Title
	research := &models.Research{
		Title:                                    &title,
		Purpose:                                  body.Purpose,
		Category:                                 body.Category,
		Funding_Scheme:                           body.FundingScheme,
		Fund_Amount:                              body.FundAmount,
		Location:                                 body.Location,
		Address:                                  body.Address,
		RPJMD:                                    body.RPJMD,
		Promised_Output:                          body.PromisedOutput,
		Research_Significance:                    body.ResearchSignificance,
		Research_Time_Range:                      body.ResearchTimeRange,
		Partnership_Target_Selection:             body.PartnershipTargetSelection,
		Partner_Name:                             body.PartnerName,
		Estimated_Number_of_Direct_Beneficiaries: body.EstimatedDirectBeneficiaries,
		File:                                     filename,
	}
	if err := c.services.Submit(&pubID, research); err != nil {
		_ = utils.FileDeletion(filename)
		return utils.HandleError(ctx, "Failed to Submit Research Proposal", err)
	}
	status := string(research.Status)
	return utils.CreatedWithData(ctx, "Proposal Submitted!", models.SubmissionResponse{
		PublicID:            research.PublicID,
		ResearchCode:        *research.ResearchCode,
		Title:               title,
		Slug:                research.Slug,
		Category:            research.Category,
		Location:            research.Location,
		Research_Time_Range: research.Research_Time_Range,
		Status:              &status,
		FundingScheme:       research.Funding_Scheme,
		FundAmount:          research.Fund_Amount,
		CreatedAt:           research.CreatedAt,
	})
}
func (c *ResearchController) MyResearchSubmission(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userInternalID := int64(claims["user_id"].(float64))
	var response []models.SubmissionResponse
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "5"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	submission, total, err := c.services.UserResearches(uint(userInternalID), filter, sort, offset, limit)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch Data", err.Error())
	}

	for _, n := range submission {
		response = append(response, models.SubmissionResponse{
			PublicID:            n.PublicID,
			ResearchCode:        *n.ResearchCode,
			Slug:                n.Slug,
			Title:               *n.Title,
			Category:            n.Category,
			Location:            n.Location,
			Research_Time_Range: n.Research_Time_Range,
			File:                &n.File,
			Status:              (*string)(&n.Status),
			UploaderPublicID:    n.Uploader.PublicID,
			UploaderName:        *n.Uploader.Name,
			FundingScheme:       n.Funding_Scheme,
			FundAmount:          n.Fund_Amount,
			CreatedAt:           n.CreatedAt,
		})
	}
	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}
	return utils.SuccessPaginated(ctx, "Submission Fetched", response, meta)
}

func (c *ResearchController) FetchResearchBySlug(ctx *fiber.Ctx) error {
	ResearchSlug := ctx.Params("slug")
	if ResearchSlug == "" {
		return utils.NotFound(ctx, "Submission Not Found", "")
	}
	research, err := c.services.ResearchSlug(ResearchSlug)
	if err != nil {
		return utils.NotFound(ctx, "Research not found", err.Error())
	}
	Response := models.Research{
		PublicID:                                 research.PublicID,
		Title:                                    research.Title,
		ResearchCode:                             research.ResearchCode,
		Slug:                                     research.Slug,
		Purpose:                                  research.Purpose,
		Category:                                 research.Category,
		Funding_Scheme:                           research.Funding_Scheme,
		Fund_Amount:                              research.Fund_Amount,
		Location:                                 research.Location,
		Address:                                  research.Address,
		RPJMD:                                    research.RPJMD,
		Promised_Output:                          research.Promised_Output,
		Research_Significance:                    research.Research_Significance,
		Research_Time_Range:                      research.Research_Time_Range,
		Partnership_Target_Selection:             research.Partnership_Target_Selection,
		Partner_Name:                             research.Partner_Name,
		Estimated_Number_of_Direct_Beneficiaries: research.Estimated_Number_of_Direct_Beneficiaries,
		File:                                     research.File,
		Status:                                   research.Status,
		CreatedAt:                                research.CreatedAt,
	}
	return utils.Success(ctx, "Research Fetched", Response)
}

func (c *ResearchController) Revision(ctx *fiber.Ctx) error {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "Missing token")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "Invalid token")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "Invalid user_id claim")
	}
	userID := int64(userIDFloat)
	publicID := ctx.Params("id")
	if publicID == "" {
		return utils.NotFound(ctx, "Submission Not Found", "Missing submission ID")
	}
	var req models.UpdateSubmission
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadReq(ctx, "Failed Request Parsing", err.Error())
	}
	if !isUUID(publicID) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	if err := c.services.Revision(userID, publicID, req); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return utils.NotFound(ctx, "Revision Failed", "submission not found")
		}
		return utils.HandleError(ctx, "Revision Failed", err)
	}

	return utils.Success(ctx, "Revision Submitted", nil)
}

//Groups

// JoinGroup: POST /api/v1/research/group/join — bergabung dengan kode (hanya akun peneliti).
func (r *ResearchController) JoinGroup(ctx *fiber.Ctx) error {
	return NewGroupController().Join(ctx)
}

// List: GET /api/v1/research/group/list — kelompok riset milik pengguna (ketua, anggota, pembimbing; admin: semua).
func (r *ResearchController) List(ctx *fiber.Ctx) error {
	return NewGroupController().Mine(ctx)
}

// GetBySlug: GET /api/v1/research/group/:slug — detail kelompok berdasarkan slug atau public_id.
func (r *ResearchController) GetBySlug(ctx *fiber.Ctx) error {
	return NewGroupController().BySlug(ctx)
}

// SyncGroups: POST /api/v1/admin/research/groups/sync — buat kelompok untuk riset
// yang sudah disetujui tetapi belum memiliki kelompok.
func (r *ResearchController) SyncGroups(ctx *fiber.Ctx) error {
	created, err := r.GService.BackfillMissingGroups()
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to sync groups", err.Error())
	}
	return utils.Success(ctx, "Groups synced", fiber.Map{"created": created})
}
func (c *ResearchController) GetByPublicID(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if publicID == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}
	group, err := c.GService.FetchByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Group Not Found", err.Error())
	}
	return utils.Success(ctx, "Group fetched", group)
}

func (c *ResearchController) GetMyGroups(ctx *fiber.Ctx) error {
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

	groups, total, err := c.GService.ListMyGroups(userID, filter, sort, offset, limit)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}
	return utils.SuccessPaginated(ctx, "My groups fetched", groups, meta)
}

func (c *ResearchController) GetMyGroup(ctx *fiber.Ctx) error {
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

	publicID := ctx.Params("id")
	if publicID == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}

	group, err := c.GService.GetMyGroup(userID, publicID)
	if err != nil {
		return utils.NotFound(ctx, "Group Not Found", err.Error())
	}
	return utils.Success(ctx, "Group fetched", group)
}

// Public Group Display
func (c *ResearchController) ListPublic(ctx *fiber.Ctx) error {
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

	groups, total, err := c.GService.ListPublic(filter, sort, offset, limit)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}
	return utils.SuccessPaginated(ctx, "Groups fetched", groups, meta)
}

func (c *ResearchController) GetPublicBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}
	group, err := c.GService.GetPublicBySlug(slug)
	if err != nil {
		return utils.NotFound(ctx, "Group Not Found", err.Error())
	}
	return utils.Success(ctx, "Group fetched", group)
}

func (c *ResearchController) GetPublicByPublicID(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if publicID == "" {
		return utils.NotFound(ctx, "Group Not Found", "")
	}
	group, err := c.GService.GetPublicByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Group Not Found", err.Error())
	}
	return utils.Success(ctx, "Group fetched", group)
}
