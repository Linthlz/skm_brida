package controllers

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"gorm.io/gorm"
)

// GroupController: kelompok riset (ketua, anggota, pembimbing BRIDA).
type GroupController struct{}

func NewGroupController() *GroupController { return &GroupController{} }

type researchSummary struct {
	PublicID     uuid.UUID `json:"public_id"`
	ResearchCode string    `json:"research_code"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Status       string    `json:"status"`
	Category     string    `json:"category"`
	Location     string    `json:"location"`
}

type groupSummary struct {
	PublicID    uuid.UUID        `json:"public_id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	MyRole      string           `json:"my_role"`
	MemberCount int              `json:"member_count"`
	Leader      *personView      `json:"leader,omitempty"`
	Supervisor  *personView      `json:"supervisor,omitempty"`
	Research    *researchSummary `json:"research,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type groupMemberView struct {
	PublicID     uuid.UUID `json:"public_id"`
	UserPublicID uuid.UUID `json:"user_public_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Institution  string    `json:"institution"`
	Position     string    `json:"position"`
	Role         string    `json:"role"`
	JoinedAt     time.Time `json:"joined_at"`
}

type supervisorView struct {
	personView
	Email string `json:"email"`
}

func nameOf(u *models.User) string {
	if u == nil || u.Name == nil {
		return ""
	}
	return *u.Name
}

func activeMembers(groupIDs []int64) (map[int64][]models.ResearchGroupMember, error) {
	var rows []models.ResearchGroupMember
	err := config.DB.Preload("User").
		Where("group_internal_id IN ? AND left_at IS NULL", groupIDs).
		Order("CASE role WHEN 'leader' THEN 0 ELSE 1 END, joined_at ASC").Find(&rows).Error
	out := map[int64][]models.ResearchGroupMember{}
	for _, m := range rows {
		out[m.GroupInternalID] = append(out[m.GroupInternalID], m)
	}
	return out, err
}

func researchesByGroup(groupIDs []int64) (map[int64]*models.Research, error) {
	var rows []models.Research
	err := config.DB.Preload("Uploader").Preload("Reviewer").
		Where("group_internal_id IN ? AND deleted_at IS NULL", groupIDs).Find(&rows).Error
	out := map[int64]*models.Research{}
	for i := range rows {
		out[*rows[i].GroupInternalID] = &rows[i]
	}
	return out, err
}

func summarize(r *models.Research) *researchSummary {
	if r == nil {
		return nil
	}
	s := &researchSummary{PublicID: r.PublicID, Slug: r.Slug, Status: string(r.Status), Category: r.Category, Location: r.Location}
	if r.ResearchCode != nil {
		s.ResearchCode = *r.ResearchCode
	}
	if r.Title != nil {
		s.Title = *r.Title
	}
	return s
}

func groupLocals(ctx *fiber.Ctx) (*models.ResearchGroup, string, int64) {
	g, _ := ctx.Locals("group").(*models.ResearchGroup)
	role, _ := ctx.Locals("group_role").(string)
	uid, _ := ctx.Locals("user_id").(int64)
	return g, role, uid
}

func notifyUser(recipientID int64, subject, message string, groupPublicID uuid.UUID) {
	ref := "group"
	n := models.Notification{
		PublicID: uuid.New(), RecipientInternalID: recipientID, Subject: subject, Message: message,
		ReferenceType: &ref, ReferencePublicID: &groupPublicID,
	}
	_ = config.DB.Omit("Research").Create(&n).Error
}

// Mine: GET /api/v1/research/group/list (juga /api/v1/users/research/group)
func (c *GroupController) Mine(ctx *fiber.Ctx) error {
	userID, _, role := claimsOf(ctx)
	myRole := map[int64]string{}
	if role == "admin" {
		var ids []int64
		if err := config.DB.Model(&models.ResearchGroup{}).Order("created_at DESC").Limit(500).Pluck("group_id", &ids).Error; err != nil {
			return utils.InternalServerError(ctx, "Failed to fetch groups", err.Error())
		}
		for _, id := range ids {
			myRole[id] = "admin"
		}
	} else {
		var members []models.ResearchGroupMember
		if err := config.DB.Where("user_internal_id = ? AND left_at IS NULL", userID).Find(&members).Error; err != nil {
			return utils.InternalServerError(ctx, "Failed to fetch groups", err.Error())
		}
		for _, m := range members {
			myRole[m.GroupInternalID] = m.Role
		}
		var supervised []int64
		if err := config.DB.Model(&models.Research{}).
			Where("reviewer_internal_id = ? AND group_internal_id IS NOT NULL AND deleted_at IS NULL", userID).
			Pluck("group_internal_id", &supervised).Error; err != nil {
			return utils.InternalServerError(ctx, "Failed to fetch groups", err.Error())
		}
		for _, id := range supervised {
			if _, ok := myRole[id]; !ok {
				myRole[id] = "pembimbing"
			}
		}
	}
	out := []groupSummary{}
	if len(myRole) == 0 {
		return utils.Success(ctx, "Groups fetched", out)
	}
	ids := make([]int64, 0, len(myRole))
	for id := range myRole {
		ids = append(ids, id)
	}
	var groups []models.ResearchGroup
	if err := config.DB.Preload("Leader").Where("group_id IN ?", ids).Find(&groups).Error; err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch groups", err.Error())
	}
	members, err := activeMembers(ids)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch groups", err.Error())
	}
	research, err := researchesByGroup(ids)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch groups", err.Error())
	}
	for _, g := range groups {
		r := research[g.InternalID]
		item := groupSummary{
			PublicID: g.PublicID, Name: g.Name, MyRole: myRole[g.InternalID],
			MemberCount: len(members[g.InternalID]), Leader: person(g.Leader), Research: summarize(r), CreatedAt: g.CreatedAt,
		}
		if r != nil {
			item.Supervisor = person(r.Reviewer)
		}
		if g.Slug != nil {
			item.Slug = *g.Slug
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return utils.Success(ctx, "Groups fetched", out)
}

// Detail: GET /api/v1/users/research/group/:id (anggota, ketua, pembimbing, admin)
func (c *GroupController) Detail(ctx *fiber.Ctx) error {
	g, myRole, _ := groupLocals(ctx)
	members, err := activeMembers([]int64{g.InternalID})
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch group", err.Error())
	}
	research, err := researchesByGroup([]int64{g.InternalID})
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch group", err.Error())
	}
	mv := make([]groupMemberView, 0, len(members[g.InternalID]))
	for _, m := range members[g.InternalID] {
		if m.User == nil {
			continue
		}
		mv = append(mv, groupMemberView{
			PublicID: m.PublicID, UserPublicID: m.User.PublicID, Name: nameOf(m.User), Email: m.User.Email,
			Institution: m.User.Institution, Position: m.User.Position, Role: m.Role, JoinedAt: m.JoinedAt,
		})
	}
	group := fiber.Map{"public_id": g.PublicID, "slug": g.Slug, "name": g.Name, "description": g.Description, "created_at": g.CreatedAt}
	if myRole == "leader" || myRole == "admin" {
		group["join_code"] = g.JoinCode
	}
	resp := fiber.Map{"group": group, "my_role": myRole, "members": mv}
	if r := research[g.InternalID]; r != nil {
		resp["research"] = toResearchView(r, true)
		if p := person(r.Reviewer); p != nil {
			resp["supervisor"] = supervisorView{personView: *p, Email: r.Reviewer.Email}
		}
	}
	return utils.Success(ctx, "Group fetched", resp)
}

// Join: POST /api/v1/research/group/join — hanya akun peneliti.
func (c *GroupController) Join(ctx *fiber.Ctx) error {
	userID, _, role := claimsOf(ctx)
	if role != "researcher" {
		return utils.Forbidden(ctx, "Cannot join group", "only researcher accounts can join a research group")
	}
	var body struct {
		JoinCode string `json:"join_code"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	code := strings.ToUpper(strings.TrimSpace(body.JoinCode))
	if code == "" || len(code) > 20 {
		return utils.BadReq(ctx, "Cannot join group", "join_code is required")
	}
	var g models.ResearchGroup
	if err := config.DB.Where("join_code = ?", code).First(&g).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound(ctx, "Cannot join group", "join code not found")
		}
		return utils.InternalServerError(ctx, "Cannot join group", err.Error())
	}
	if g.LeaderInternalID == userID {
		return utils.Conflict(ctx, "Cannot join group", "you are the leader of this group")
	}
	var existing models.ResearchGroupMember
	err := config.DB.Unscoped().Where("group_internal_id = ? AND user_internal_id = ?", g.InternalID, userID).First(&existing).Error
	now := time.Now()
	switch {
	case err == nil && existing.LeftAt == nil && !existing.DeletedAt.Valid:
		return utils.Conflict(ctx, "Cannot join group", "you are already a member of this group")
	case err == nil:
		if err := config.DB.Unscoped().Model(&models.ResearchGroupMember{}).Where("member_id = ?", existing.InternalID).
			Updates(map[string]interface{}{"left_at": nil, "deleted_at": nil, "joined_at": now, "role": "member"}).Error; err != nil {
			return utils.InternalServerError(ctx, "Cannot join group", err.Error())
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		m := models.ResearchGroupMember{PublicID: uuid.New(), GroupInternalID: g.InternalID, UserInternalID: userID, Role: "member", JoinedAt: now}
		if err := config.DB.Omit("User").Create(&m).Error; err != nil {
			return utils.InternalServerError(ctx, "Cannot join group", err.Error())
		}
	default:
		return utils.InternalServerError(ctx, "Cannot join group", err.Error())
	}
	var joiner models.User
	config.DB.First(&joiner, "internal_id = ?", userID)
	notifyUser(g.LeaderInternalID, "Anggota baru bergabung",
		fmt.Sprintf("%s bergabung ke kelompok riset \"%s\".", nameOf(&joiner), g.Name), g.PublicID)
	return utils.Success(ctx, "Joined group", fiber.Map{"public_id": g.PublicID, "name": g.Name})
}

type researchEditRequest struct {
	Purpose                    *string `json:"purpose"`
	ResearchSignificance       *string `json:"research_significance"`
	PromisedOutput             *string `json:"promised_output"`
	ResearchTimeRange          *string `json:"research_time_range"`
	Location                   *string `json:"location"`
	Address                    *string `json:"address"`
	PartnershipTargetSelection *string `json:"partnership_target_selection"`
	PartnerName                *string `json:"partner_name"`
	Beneficiaries              *string `json:"estimated_number_of_direct_beneficiaries"`
}

// UpdateResearch: PATCH /api/v1/research/group/:id/research — hanya ketua peneliti.
// Judul, bidang, skema, dana, dan RPJMD terkunci setelah disetujui BRIDA.
func (c *GroupController) UpdateResearch(ctx *fiber.Ctx) error {
	g, myRole, _ := groupLocals(ctx)
	if myRole != "leader" {
		return utils.Forbidden(ctx, "Cannot update research", "only the lead researcher can edit this research")
	}
	var body researchEditRequest
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	fields := []struct {
		col string
		val *string
		max int
		req bool
	}{
		{"purpose", body.Purpose, 20000, true},
		{"research_significance", body.ResearchSignificance, 20000, false},
		{"promised_output", body.PromisedOutput, 20000, true},
		{"research_time_range", body.ResearchTimeRange, 255, false},
		{"location", body.Location, 255, false},
		{"address", body.Address, 2000, false},
		{"partnership_target_selection", body.PartnershipTargetSelection, 2000, false},
		{"partner_name", body.PartnerName, 255, false},
		{"estimated_number_of_direct_beneficiaries", body.Beneficiaries, 255, false},
	}
	updates := map[string]interface{}{}
	for _, f := range fields {
		if f.val == nil {
			continue
		}
		v := strings.TrimSpace(*f.val)
		if f.req && v == "" {
			return utils.BadReq(ctx, "Cannot update research", f.col+" cannot be empty")
		}
		if utf8.RuneCountInString(v) > f.max {
			return utils.BadReq(ctx, "Cannot update research", fmt.Sprintf("%s must be at most %d characters", f.col, f.max))
		}
		updates[f.col] = v
	}
	if len(updates) == 0 {
		return utils.BadReq(ctx, "Cannot update research", "no editable field was sent")
	}
	res := config.DB.Model(&models.Research{}).
		Where("group_internal_id = ? AND deleted_at IS NULL AND status IN ?", g.InternalID, []string{"on-going", "approved"}).
		Updates(updates)
	if res.Error != nil {
		return utils.InternalServerError(ctx, "Cannot update research", res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return utils.NotFound(ctx, "Cannot update research", "no running research for this group")
	}
	return c.Detail(ctx)
}

// RemoveMember: DELETE /api/v1/research/group/:id/members/:memberID — ketua atau admin.
func (c *GroupController) RemoveMember(ctx *fiber.Ctx) error {
	g, myRole, _ := groupLocals(ctx)
	if myRole != "leader" && myRole != "admin" {
		return utils.Forbidden(ctx, "Cannot remove member", "only the lead researcher can remove members")
	}
	memberID := ctx.Params("memberID")
	if !isUUID(memberID) {
		return utils.BadReq(ctx, "Invalid ID", "memberID must be a valid UUID")
	}
	var m models.ResearchGroupMember
	if err := config.DB.Where("public_id = ? AND group_internal_id = ? AND left_at IS NULL", memberID, g.InternalID).First(&m).Error; err != nil {
		return utils.NotFound(ctx, "Cannot remove member", "member not found")
	}
	if m.Role == "leader" {
		return utils.Conflict(ctx, "Cannot remove member", "the lead researcher cannot be removed")
	}
	if err := config.DB.Model(&models.ResearchGroupMember{}).Where("member_id = ?", m.InternalID).Update("left_at", time.Now()).Error; err != nil {
		return utils.InternalServerError(ctx, "Cannot remove member", err.Error())
	}
	notifyUser(m.UserInternalID, "Dikeluarkan dari kelompok riset",
		fmt.Sprintf("Anda tidak lagi menjadi anggota kelompok riset \"%s\".", g.Name), g.PublicID)
	return utils.Success(ctx, "Member removed", nil)
}

// Leave: POST /api/v1/research/group/:id/leave — anggota (bukan ketua).
func (c *GroupController) Leave(ctx *fiber.Ctx) error {
	g, myRole, uid := groupLocals(ctx)
	if myRole != "member" {
		return utils.Forbidden(ctx, "Cannot leave group", "only members can leave; the lead researcher stays with the research")
	}
	if err := config.DB.Model(&models.ResearchGroupMember{}).
		Where("group_internal_id = ? AND user_internal_id = ? AND left_at IS NULL", g.InternalID, uid).
		Update("left_at", time.Now()).Error; err != nil {
		return utils.InternalServerError(ctx, "Cannot leave group", err.Error())
	}
	return utils.Success(ctx, "Left group", nil)
}

// RegenerateCode: POST /api/v1/research/group/:id/join-code — ketua atau admin.
func (c *GroupController) RegenerateCode(ctx *fiber.Ctx) error {
	g, myRole, _ := groupLocals(ctx)
	if myRole != "leader" && myRole != "admin" {
		return utils.Forbidden(ctx, "Cannot change join code", "only the lead researcher can change the join code")
	}
	code := utils.GenerateJoinCode(g.Name)
	if err := config.DB.Model(&models.ResearchGroup{}).Where("group_id = ?", g.InternalID).Update("join_code", code).Error; err != nil {
		return utils.InternalServerError(ctx, "Cannot change join code", err.Error())
	}
	return utils.Success(ctx, "Join code changed", fiber.Map{"join_code": code})
}

// SetSupervisor: PATCH /api/v1/research/group/:id/supervisor — admin menetapkan pembimbing BRIDA.
func (c *GroupController) SetSupervisor(ctx *fiber.Ctx) error {
	g, myRole, _ := groupLocals(ctx)
	if myRole != "admin" {
		return utils.Forbidden(ctx, "Cannot set supervisor", "only administrators can assign a supervisor")
	}
	var body struct {
		ReviewerPublicID string `json:"reviewer_public_id"`
	}
	if err := ctx.BodyParser(&body); err != nil || !isUUID(body.ReviewerPublicID) {
		return utils.BadReq(ctx, "Cannot set supervisor", "reviewer_public_id must be a valid UUID")
	}
	var reviewer models.User
	if err := config.DB.Where("public_id = ?", body.ReviewerPublicID).First(&reviewer).Error; err != nil {
		return utils.NotFound(ctx, "Cannot set supervisor", "reviewer not found")
	}
	if reviewer.Role != "reviewer" {
		return utils.BadReq(ctx, "Cannot set supervisor", "the selected user is not a reviewer")
	}
	res := config.DB.Model(&models.Research{}).Where("group_internal_id = ? AND deleted_at IS NULL", g.InternalID).
		Update("reviewer_internal_id", reviewer.InternalID)
	if res.Error != nil {
		return utils.InternalServerError(ctx, "Cannot set supervisor", res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return utils.NotFound(ctx, "Cannot set supervisor", "no research for this group")
	}
	notifyUser(reviewer.InternalID, "Ditugaskan sebagai pembimbing",
		fmt.Sprintf("Anda ditugaskan sebagai pembimbing BRIDA untuk kelompok riset \"%s\".", g.Name), g.PublicID)
	return c.Detail(ctx)
}

// BySlug: GET /api/v1/research/group/:slug — detail kelompok berdasarkan slug (atau public_id),
// dengan pemeriksaan akses yang sama seperti middleware GroupMemberOnly.
func (c *GroupController) BySlug(ctx *fiber.Ctx) error {
	userID, _, role := claimsOf(ctx)
	key := strings.TrimSpace(ctx.Params("slug"))
	if key == "" {
		return utils.BadReq(ctx, "Invalid group", "slug is required")
	}
	q := config.DB.Preload("Leader")
	if isUUID(key) {
		q = q.Where("public_id = ?", key)
	} else {
		q = q.Where("slug = ?", key)
	}
	var g models.ResearchGroup
	if err := q.First(&g).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound(ctx, "Group Not Found", "group not found")
		}
		return utils.InternalServerError(ctx, "Failed to fetch group", err.Error())
	}
	repo := repositories.NewResearchGroupRepository()
	groupRole := ""
	switch {
	case strings.EqualFold(role, "admin"):
		groupRole = "admin"
	case g.LeaderInternalID == userID:
		groupRole = "leader"
	default:
		isMember, err := repo.IsMember(g.InternalID, userID)
		if err != nil {
			return utils.InternalServerError(ctx, "Failed to check membership", err.Error())
		}
		if isMember {
			groupRole = "member"
		} else if isSupervisor, err := repo.IsSupervisor(g.InternalID, userID); err != nil {
			return utils.InternalServerError(ctx, "Failed to check membership", err.Error())
		} else if isSupervisor {
			groupRole = "pembimbing"
		}
	}
	if groupRole == "" {
		return utils.Forbidden(ctx, "Forbidden", "You are not a member of this group")
	}
	ctx.Locals("group", &g)
	ctx.Locals("user_id", userID)
	ctx.Locals("group_role", groupRole)
	return c.Detail(ctx)
}
