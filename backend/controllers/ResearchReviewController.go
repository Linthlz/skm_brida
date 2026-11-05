package controllers

import (
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type ResearchReviewController struct {
	services        services.SubmissionCommentService
	researchService services.ResearchService
}

func NewResearchReviewController(services services.SubmissionCommentService, researchService services.ResearchService) *ResearchReviewController {
	return &ResearchReviewController{services: services, researchService: researchService}
}

func (c *ResearchReviewController) AssignReviewer(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	FreviewerID, ok := claims["user_id"].(float64)
	reviewerID := int64(FreviewerID)
	SubmissionPublicID := ctx.Params("id")
	if SubmissionPublicID == "" {
		return utils.NotFound(ctx, "Submission Not Found", "")
	}
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	if err := c.services.AssignReviewer(reviewerID, SubmissionPublicID); err != nil {
		return utils.BadReq(ctx, "Bad Request", err.Error())
	}
	return utils.Success(ctx, "Reviewer Assigned", "")

}
func (c *ResearchReviewController) ApproveSubmission(ctx *fiber.Ctx) error {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	FreviewerID, ok := claims["user_id"].(float64)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "")
	}
	reviewerID := int64(FreviewerID)
	SubmissionPublicID := ctx.Params("id")
	if SubmissionPublicID == "" {
		return utils.NotFound(ctx, "Submission Not Found", "")
	}
	//if you want to leave a message
	var body struct {
		Message string `json:"message"`
	}
	_ = ctx.BodyParser(&body)
	// Default message empty
	if body.Message == "" {
		body.Message = "Approved."
	}
	if err := c.services.ApproveSubmission(reviewerID, SubmissionPublicID, body.Message); err != nil {
		return utils.BadReq(ctx, "Bad Request", err.Error())
	}
	return utils.Success(ctx, "Research is now on going!", "")
}

func (c *ResearchReviewController) AddComment(ctx *fiber.Ctx) error {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "Missing token")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "Invalid token")
	}
	reviewerIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Access Forbidden", "Invalid user_id claim")
	}
	reviewerID := int64(reviewerIDFloat)
	submissionPublicID := ctx.Params("id")
	if submissionPublicID == "" {
		return utils.NotFound(ctx, "Submission Not Found", "Missing submission ID")
	}
	var body struct {
		Decision string `json:"decision"`
		Message  string `json:"message"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	normalized := strings.ToLower(strings.TrimSpace(body.Decision))
	decision := models.SubmissionDecision(normalized)
	comment, err := c.services.AddReview(reviewerID, submissionPublicID, decision, body.Message)
	if err != nil {
		return utils.BadReq(ctx, "Unable to submit review", err.Error())
	}

	return utils.Success(ctx, "Review submitted", comment)
}
func (c *ResearchReviewController) Researches(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "5"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 5
	}
	offset := (page - 1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")
	submission, total, err := c.researchService.Researches(filter, sort, offset, limit)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch Data", err.Error())
	}
	response := make([]models.SubmissionResponse, 0, len(submission))
	for _, n := range submission {
		item := models.SubmissionResponse{
			PublicID:            n.PublicID,
			Slug:                n.Slug,
			Category:            n.Category,
			Location:            n.Location,
			Research_Time_Range: n.Research_Time_Range,
			File:                &n.File,
			Status:              (*string)(&n.Status),
		}
		if n.ResearchCode != nil {
			item.ResearchCode = *n.ResearchCode
		}
		if n.Title != nil {
			item.Title = *n.Title
		}
		if n.Uploader != nil {
			item.UploaderPublicID = n.Uploader.PublicID
			if n.Uploader.Name != nil {
				item.UploaderName = *n.Uploader.Name
			}
		}
		response = append(response, item)
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
