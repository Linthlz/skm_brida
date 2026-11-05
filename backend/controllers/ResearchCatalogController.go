package controllers

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

// Statuses visible in the public catalog.
var publicResearchStatuses = []string{string(models.SubmissionOnGoing), string(models.SubmissionApproved)}

// Roles allowed to read any submission in full.
var staffRoles = map[string]bool{"admin": true, "reviewer": true, "opd": true, "pegawai-brida": true}

type ResearchCatalogController struct{}

func NewResearchCatalogController() *ResearchCatalogController { return &ResearchCatalogController{} }

type personView struct {
	PublicID    uuid.UUID `json:"public_id"`
	Name        string    `json:"name"`
	Institution string    `json:"institution"`
	Position    string    `json:"position"`
}

type researchView struct {
	PublicID                   uuid.UUID   `json:"public_id"`
	ResearchCode               string      `json:"research_code"`
	Slug                       string      `json:"slug"`
	Title                      string      `json:"title"`
	Purpose                    string      `json:"purpose"`
	Category                   string      `json:"category"`
	FundingScheme              string      `json:"funding_scheme"`
	FundAmount                 string      `json:"fund_amount"`
	Location                   string      `json:"location"`
	Address                    string      `json:"address"`
	RPJMD                      string      `json:"rpjmd"`
	PromisedOutput             string      `json:"promised_output"`
	ResearchSignificance       string      `json:"research_significance"`
	ResearchTimeRange          string      `json:"research_time_range"`
	PartnershipTargetSelection string      `json:"partnership_target_selection"`
	PartnerName                string      `json:"partner_name"`
	Beneficiaries              string      `json:"estimated_number_of_direct_beneficiaries"`
	Status                     string      `json:"status"`
	ProposalFile               string      `json:"research_proposal_file,omitempty"`
	Uploader                   *personView `json:"uploader,omitempty"`
	Reviewer                   *personView `json:"reviewer,omitempty"`
	GroupPublicID              string      `json:"group_public_id,omitempty"`
	CreatedAt                  time.Time   `json:"created_at"`
	UpdatedAt                  time.Time   `json:"updated_at"`
}

type commentView struct {
	PublicID       uuid.UUID `json:"public_id"`
	AuthorName     string    `json:"author_name"`
	AuthorPosition string    `json:"author_position"`
	Message        string    `json:"message"`
	Decision       string    `json:"decision,omitempty"`
	Round          int       `json:"round,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type memberView struct {
	Name        string `json:"name"`
	Institution string `json:"institution"`
	Role        string `json:"role"`
}

type groupView struct {
	PublicID uuid.UUID    `json:"public_id"`
	Name     string       `json:"name"`
	JoinCode string       `json:"join_code,omitempty"`
	Members  []memberView `json:"members"`
}

func person(u *models.User) *personView {
	if u == nil || u.InternalID == 0 {
		return nil
	}
	name := ""
	if u.Name != nil {
		name = *u.Name
	}
	return &personView{PublicID: u.PublicID, Name: name, Institution: u.Institution, Position: u.Position}
}

func toResearchView(r *models.Research, withFile bool) researchView {
	v := researchView{
		PublicID: r.PublicID, Slug: r.Slug, Purpose: r.Purpose, Category: r.Category,
		FundingScheme: r.Funding_Scheme, FundAmount: r.Fund_Amount, Location: r.Location, Address: r.Address,
		RPJMD: r.RPJMD, PromisedOutput: r.Promised_Output, ResearchSignificance: r.Research_Significance,
		ResearchTimeRange: r.Research_Time_Range, PartnershipTargetSelection: r.Partnership_Target_Selection,
		PartnerName: r.Partner_Name, Beneficiaries: r.Estimated_Number_of_Direct_Beneficiaries,
		Status: string(r.Status), Uploader: person(r.Uploader), Reviewer: person(r.Reviewer),
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
	if r.ResearchCode != nil {
		v.ResearchCode = *r.ResearchCode
	}
	if r.Title != nil {
		v.Title = *r.Title
	}
	if r.Group != nil && r.Group.InternalID != 0 {
		v.GroupPublicID = r.Group.PublicID.String()
	}
	if withFile && r.File != "" {
		v.ProposalFile = "/" + strings.TrimPrefix(strings.TrimPrefix(r.File, "./"), "public/")
	}
	return v
}

func groupFor(r *models.Research, withJoinCode bool) (*groupView, error) {
	if r.GroupInternalID == nil {
		return nil, nil
	}
	g, err := repositories.FetchGroupWithMembers(*r.GroupInternalID)
	if err != nil || g == nil {
		return nil, err
	}
	view := &groupView{PublicID: g.PublicID, Name: g.Name, Members: []memberView{}}
	if withJoinCode {
		view.JoinCode = g.JoinCode
	}
	for _, m := range g.Members {
		if m.User == nil {
			continue
		}
		name := ""
		if m.User.Name != nil {
			name = *m.User.Name
		}
		view.Members = append(view.Members, memberView{Name: name, Institution: m.User.Institution, Role: m.Role})
	}
	return view, nil
}

func paging(ctx *fiber.Ctx, defLimit int) (page, limit, offset int) {
	page, _ = strconv.Atoi(ctx.Query("page", "1"))
	limit, _ = strconv.Atoi(ctx.Query("limit", strconv.Itoa(defLimit)))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = defLimit
	}
	return page, limit, (page - 1) * limit
}

func metaFor(page, limit int, total int64, filter, sort string) utils.PaginationMeta {
	return utils.PaginationMeta{Page: page, Limit: limit, Total: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))), Filter: filter, Sort: sort}
}

func claimsOf(ctx *fiber.Ctx) (userID int64, pubID, role string) {
	token, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return 0, "", ""
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	if id, ok := claims["user_id"].(float64); ok {
		userID = int64(id)
	}
	pubID, _ = claims["pub_id"].(string)
	role, _ = claims["role"].(string)
	return userID, pubID, role
}

func isGroupMember(groupID *int64, userID int64) bool {
	if groupID == nil {
		return false
	}
	var count int64
	config.DB.Table("research_group_members").
		Where("group_internal_id = ? AND user_internal_id = ? AND left_at IS NULL AND deleted_at IS NULL", *groupID, userID).Count(&count)
	return count > 0
}

// PublicList: GET /api/v2/researches
func (c *ResearchCatalogController) PublicList(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 12)
	q := repositories.ResearchQuery{
		Statuses: publicResearchStatuses,
		Filter:   strings.TrimSpace(ctx.Query("filter")),
		Category: strings.TrimSpace(ctx.Query("category")),
		Location: strings.TrimSpace(ctx.Query("location")),
		Sort:     ctx.Query("sort"),
		Offset:   offset, Limit: limit,
	}
	if len(q.Filter) > 100 {
		return utils.BadReq(ctx, "Invalid filter", "filter is too long")
	}
	list, total, err := repositories.SearchResearches(q)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch researches", err.Error())
	}
	out := make([]researchView, 0, len(list))
	for i := range list {
		out = append(out, toResearchView(&list[i], false))
	}
	return utils.SuccessPaginated(ctx, "Researches fetched", out, metaFor(page, limit, total, q.Filter, q.Sort))
}

// PublicFacets: GET /api/v2/researches/facets?by=category|location
func (c *ResearchCatalogController) PublicFacets(ctx *fiber.Ctx) error {
	by := ctx.Query("by", "category")
	rows, err := repositories.ResearchFacets(by, publicResearchStatuses)
	if err != nil {
		if by != "category" && by != "location" {
			return utils.BadReq(ctx, "Invalid facet", "by must be 'category' or 'location'")
		}
		return utils.InternalServerError(ctx, "Failed to fetch facets", err.Error())
	}
	return utils.Success(ctx, "Facets fetched", rows)
}

// PublicDetail: GET /api/v2/researches/:slug
func (c *ResearchCatalogController) PublicDetail(ctx *fiber.Ctx) error {
	r, err := repositories.FetchResearchBySlug(ctx.Params("slug"))
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch research", err.Error())
	}
	if r == nil || (r.Status != models.SubmissionOnGoing && r.Status != models.SubmissionApproved) {
		return utils.NotFound(ctx, "Research not found", "research not found")
	}
	group, err := groupFor(r, false)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch research", err.Error())
	}
	return utils.Success(ctx, "Research fetched", fiber.Map{"research": toResearchView(r, false), "group": group})
}

// PublicDocumentations: GET /api/v2/documentations
func (c *ResearchCatalogController) PublicDocumentations(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 12)
	docs, total, err := repositories.PublicDocumentations(repositories.DocumentationQuery{
		ResearchSlug: strings.TrimSpace(ctx.Query("research")),
		Category:     strings.TrimSpace(ctx.Query("category")),
		Offset:       offset, Limit: limit,
	})
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch documentations", err.Error())
	}
	return utils.SuccessPaginated(ctx, "Documentations fetched", docs, metaFor(page, limit, total, "", ""))
}

// StaffList: GET /api/v1/reviewer/researches?status=pending,under-review
func (c *ResearchCatalogController) StaffList(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 20)
	var statuses []string
	for _, s := range strings.Split(ctx.Query("status"), ",") {
		if s = strings.TrimSpace(s); s != "" {
			statuses = append(statuses, s)
		}
	}
	q := repositories.ResearchQuery{Statuses: statuses, Filter: strings.TrimSpace(ctx.Query("filter")), Sort: ctx.Query("sort"), Offset: offset, Limit: limit}
	list, total, err := repositories.SearchResearches(q)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch researches", err.Error())
	}
	out := make([]researchView, 0, len(list))
	for i := range list {
		out = append(out, toResearchView(&list[i], true))
	}
	return utils.SuccessPaginated(ctx, "Researches fetched", out, metaFor(page, limit, total, q.Filter, q.Sort))
}

// Detail: GET /api/v1/research/:slug — owner, group member or staff only.
func (c *ResearchCatalogController) Detail(ctx *fiber.Ctx) error {
	userID, _, role := claimsOf(ctx)
	r, err := repositories.FetchResearchBySlug(ctx.Params("slug"))
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch research", err.Error())
	}
	if r == nil {
		return utils.NotFound(ctx, "Research not found", "research not found")
	}
	isOwner := r.UploaderInternalID != nil && *r.UploaderInternalID == userID
	isMember := isGroupMember(r.GroupInternalID, userID)
	if !isOwner && !isMember && !staffRoles[role] {
		return utils.Forbidden(ctx, "Forbidden", "you cannot view this submission")
	}
	comments, err := repositories.FetchResearchComments(*r.InternalID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch research", err.Error())
	}
	cv := make([]commentView, 0, len(comments))
	for _, cm := range comments {
		item := commentView{PublicID: cm.PublicID, Message: cm.Message, CreatedAt: cm.CreatedAt, AuthorPosition: cm.Author.Position}
		if cm.Author.Name != nil {
			item.AuthorName = *cm.Author.Name
		}
		if cm.Decision != nil {
			item.Decision = string(*cm.Decision)
		}
		if cm.Round != nil {
			item.Round = *cm.Round
		}
		cv = append(cv, item)
	}
	group, err := groupFor(r, isOwner || isMember || role == "admin")
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch research", err.Error())
	}
	return utils.Success(ctx, "Research fetched", fiber.Map{
		"research": toResearchView(r, true),
		"comments": cv,
		"group":    group,
	})
}
