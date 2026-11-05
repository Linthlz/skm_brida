package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

var schemeCodes = map[string]bool{"hibah-daerah": true, "insentif": true, "kolaboratif": true, "mandiri": true, "brin": true}
var fundingStatuses = map[string]bool{models.FundingOpen: true, models.FundingClosing: true, models.FundingSoon: true}

type ContentController struct{}

func NewContentController() *ContentController { return &ContentController{} }

type fundingRequest struct {
	Name         string `json:"name"`
	Organizer    string `json:"organizer"`
	SchemeCode   string `json:"scheme_code"`
	Ceiling      int64  `json:"ceiling"`
	Quota        int    `json:"quota"`
	Deadline     string `json:"deadline"`
	Status       string `json:"status"`
	TargetFields string `json:"target_fields"`
	Requirements string `json:"requirements"`
	Description  string `json:"description"`
	Website      string `json:"website"`
	WebsiteName  string `json:"website_name"`
}

func (r *fundingRequest) toModel() (*models.FundingScheme, string) {
	r.Name, r.Organizer = strings.TrimSpace(r.Name), strings.TrimSpace(r.Organizer)
	r.Website, r.WebsiteName = strings.TrimSpace(r.Website), strings.TrimSpace(r.WebsiteName)
	switch {
	case r.Name == "" || utf8.RuneCountInString(r.Name) > 255:
		return nil, "name is required (max 255 characters)"
	case r.Organizer == "" || utf8.RuneCountInString(r.Organizer) > 255:
		return nil, "organizer is required (max 255 characters)"
	case !schemeCodes[r.SchemeCode]:
		return nil, "scheme_code is not valid"
	case !fundingStatuses[r.Status]:
		return nil, "status must be open, closing or soon"
	case r.Ceiling < 0 || r.Quota < 0:
		return nil, "ceiling and quota cannot be negative"
	case utf8.RuneCountInString(r.WebsiteName) > 255 || len(r.Website) > 500:
		return nil, "website fields are too long"
	}
	if r.Website != "" {
		u, err := url.Parse(r.Website)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return nil, "website must be a valid http(s) URL"
		}
	}
	deadline, err := time.Parse("2006-01-02", r.Deadline)
	if err != nil {
		return nil, "deadline must use YYYY-MM-DD format"
	}
	return &models.FundingScheme{
		Name: r.Name, Organizer: r.Organizer, SchemeCode: r.SchemeCode, Ceiling: r.Ceiling, Quota: r.Quota,
		Deadline: deadline, Status: r.Status, TargetFields: strings.TrimSpace(r.TargetFields),
		Requirements: strings.TrimSpace(r.Requirements), Description: strings.TrimSpace(r.Description),
		Website: r.Website, WebsiteName: r.WebsiteName,
	}, ""
}

// ListFunding: GET /api/v2/funding (public) and /api/v1/funding (admin).
func (c *ContentController) ListFunding(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 50)
	status, scheme := ctx.Query("status"), ctx.Query("scheme")
	if status != "" && !fundingStatuses[status] {
		return utils.BadReq(ctx, "Invalid status", "status must be open, closing or soon")
	}
	list, total, err := repositories.ListFunding(status, scheme, offset, limit)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch funding schemes", err.Error())
	}
	return utils.SuccessPaginated(ctx, "Funding schemes fetched", list, metaFor(page, limit, total, "", ""))
}

func (c *ContentController) CreateFunding(ctx *fiber.Ctx) error {
	var body fundingRequest
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	f, msg := body.toModel()
	if msg != "" {
		return utils.BadReq(ctx, "Cannot create funding scheme", msg)
	}
	f.PublicID = uuid.New()
	if err := repositories.CreateFunding(f); err != nil {
		return utils.InternalServerError(ctx, "Cannot create funding scheme", err.Error())
	}
	return utils.CreatedWithData(ctx, "Funding scheme created", f)
}

func (c *ContentController) UpdateFunding(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !isUUID(id) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	existing, err := repositories.FetchFunding(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Cannot update funding scheme", err.Error())
	}
	if existing == nil {
		return utils.NotFound(ctx, "Cannot update funding scheme", "funding scheme not found")
	}
	var body fundingRequest
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	f, msg := body.toModel()
	if msg != "" {
		return utils.BadReq(ctx, "Cannot update funding scheme", msg)
	}
	f.InternalID, f.PublicID, f.CreatedAt = existing.InternalID, existing.PublicID, existing.CreatedAt
	if err := repositories.UpdateFunding(f); err != nil {
		return utils.InternalServerError(ctx, "Cannot update funding scheme", err.Error())
	}
	return utils.Success(ctx, "Funding scheme updated", f)
}

func (c *ContentController) DeleteFunding(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !isUUID(id) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	existing, err := repositories.FetchFunding(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Cannot delete funding scheme", err.Error())
	}
	if existing == nil {
		return utils.NotFound(ctx, "Cannot delete funding scheme", "funding scheme not found")
	}
	if err := repositories.DeleteFunding(id); err != nil {
		return utils.InternalServerError(ctx, "Cannot delete funding scheme", err.Error())
	}
	return utils.Success(ctx, "Funding scheme deleted", nil)
}

// Stats: GET /api/v2/stats
func (c *ContentController) Stats(ctx *fiber.Ctx) error {
	stats, err := repositories.PortalStatistics()
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch statistics", err.Error())
	}
	return utils.Success(ctx, "Statistics fetched", stats)
}

// GetSettings: GET /api/v2/settings
func (c *ContentController) GetSettings(ctx *fiber.Ctx) error {
	values, err := repositories.GetSettings()
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch settings", err.Error())
	}
	return utils.Success(ctx, "Settings fetched", values)
}

// SaveSettings: PUT /api/v1/admin/settings — partial update of known keys.
func (c *ContentController) SaveSettings(ctx *fiber.Ctx) error {
	var body map[string]interface{}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	values := map[string]string{}
	for k, raw := range body {
		if _, ok := repositories.SettingDefaults[k]; !ok {
			return utils.BadReq(ctx, "Cannot save settings", fmt.Sprintf("unknown setting %q", k))
		}
		var v string
		switch t := raw.(type) {
		case string:
			v = strings.TrimSpace(t)
		case bool:
			v = strconv.FormatBool(t)
		case float64:
			v = strconv.FormatInt(int64(t), 10)
		default:
			return utils.BadReq(ctx, "Cannot save settings", fmt.Sprintf("invalid value for %q", k))
		}
		if repositories.BoolSettings[k] {
			if _, err := strconv.ParseBool(v); err != nil {
				return utils.BadReq(ctx, "Cannot save settings", fmt.Sprintf("%q must be true or false", k))
			}
		} else if repositories.IntSetting(k) {
			if n, err := strconv.Atoi(v); err != nil || n < 0 {
				return utils.BadReq(ctx, "Cannot save settings", fmt.Sprintf("%q must be a non-negative number", k))
			}
		} else if v == "" && (k == "site_name" || k == "contact_email") {
			return utils.BadReq(ctx, "Cannot save settings", fmt.Sprintf("%q is required", k))
		} else if utf8.RuneCountInString(v) > 500 {
			return utils.BadReq(ctx, "Cannot save settings", fmt.Sprintf("%q is too long", k))
		}
		if k == "contact_email" && !emailPattern.MatchString(v) {
			return utils.BadReq(ctx, "Cannot save settings", "contact_email is not valid")
		}
		values[k] = v
	}
	if err := repositories.SaveSettings(values); err != nil {
		return utils.InternalServerError(ctx, "Cannot save settings", err.Error())
	}
	return c.GetSettings(ctx)
}
